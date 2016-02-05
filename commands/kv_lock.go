package commands

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type KvLockOptions struct {
	Behavior     string
	Ttl          string
	LockDelay    time.Duration
	Session      string
	cleanSession bool
}

func (k *Kv) AddLockSub(cmd *cobra.Command) {
	klo := &KvLockOptions{cleanSession: false}

	lockCmd := &cobra.Command{
		Use:   "lock <path>",
		Short: "Acquire a lock on a given path",
		Long:  "Acquire a lock on a given path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Lock(args, klo)
		},
	}

	oldLockCmd := &cobra.Command{
		Use:        "kv-lock <path>",
		Short:      "Acquire a lock on a given path",
		Long:       "Acquire a lock on a given path",
		Deprecated: "Use kv lock",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Lock(args, klo)
		},
	}

	lockCmd.Flags().StringVar(&klo.Behavior, "behavior", "release", "Lock behavior. One of 'release' or 'delete'")
	lockCmd.Flags().StringVar(&klo.Ttl, "ttl", "", "Lock time to live")
	lockCmd.Flags().DurationVar(&klo.LockDelay, "lock-delay", 15*time.Second, "Lock delay")
	lockCmd.Flags().StringVar(&klo.Session, "session", "", "Previously created session to use for lock")
	k.AddDatacenterOption(lockCmd)

	oldLockCmd.Flags().StringVar(&klo.Behavior, "behavior", "release", "Lock behavior. One of 'release' or 'delete'")
	oldLockCmd.Flags().StringVar(&klo.Ttl, "ttl", "", "Lock time to live")
	oldLockCmd.Flags().DurationVar(&klo.LockDelay, "lock-delay", 15*time.Second, "Lock delay")
	oldLockCmd.Flags().StringVar(&klo.Session, "session", "", "Previously created session to use for lock")
	k.AddDatacenterOption(oldLockCmd)

	cmd.AddCommand(lockCmd)

	k.AddCommand(oldLockCmd)
}

func (k *Kv) Lock(args []string, klo *KvLockOptions) error {
	var lockOpts *consulapi.KVPair

	// Work around a Consul API bug that ignores LockDelay == 0
	if klo.LockDelay == 0 {
		klo.LockDelay = time.Nanosecond
	}

	switch {
	case len(args) == 0:
		return fmt.Errorf("Key path must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one path allowed")
	}
	path := args[0]

	client, err := k.KV()
	if err != nil {
		return err
	}

	writeOpts := k.WriteOptions()
	queryOpts := k.QueryOptions()
	queryOpts.WaitTime = 15 * time.Second
	
	sessionClient, err := k.Session()
	if err != nil {
		return err
	}

	if klo.Session == "" {
		// Create the Consul session
		se, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
			Name:      "Session for consul-cli",
			LockDelay: klo.LockDelay,
			Behavior:  klo.Behavior,
			TTL:       klo.Ttl,
		}, writeOpts)
		if err != nil {
			return err
		}

		klo.Session = se
		klo.cleanSession = true
	}

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(klo.Ttl, klo.Session, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

WAIT:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		k.destroySession(sessionClient, klo)
		return err
	}

	locked := false

	if kv != nil && kv.Session == klo.Session {
		goto HELD
	}
	if kv != nil && kv.Session != "" {
		queryOpts.WaitIndex = meta.LastIndex
		goto WAIT
	}

	// Node doesn't already exist
	if kv == nil {
		lockOpts = &consulapi.KVPair{
			Key:     path,
			Session: klo.Session,
		}
	} else {
		lockOpts = &consulapi.KVPair{
			Key:     kv.Key,
			Flags:   kv.Flags,
			Value:   kv.Value,
			Session: klo.Session,
		}
	}

	// Try to acquire the lock
	locked, _, err = client.Acquire(lockOpts, nil)
	if err != nil {
		k.destroySession(sessionClient, klo)
		return err
	}

	if !locked {
		select {
		case <-time.After(5 * time.Second):
			goto WAIT
		}
	}

HELD:
	fmt.Fprintln(k.Out, klo.Session)

	return nil
}

// Destroy the session on error. Only performed when
// klo.cleanSession == true
//
func (k *Kv) destroySession(s *consulapi.Session, klo *KvLockOptions) error {
	if klo.cleanSession {
		writeOpts := k.WriteOptions()
		_, err := s.Destroy(klo.Session, writeOpts)
		if err != nil {
			return fmt.Errorf("Session not destroyed: %s", klo.Session)
		}
	}

	return nil
}
