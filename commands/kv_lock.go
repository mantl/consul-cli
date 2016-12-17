package commands

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Lock functions

func newKvLockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock <path>",
		Short: "Acquire a lock on a given path",
		Long:  "Acquire a lock on a given path",
		RunE:  kvLock,
	}
	cmd.Flags().String("behavior", "release", "Lock behavior. One of 'release' or 'delete'")
	cmd.Flags().String("ttl", "", "Lock time to live")
	cmd.Flags().Duration("lock-delay", 15*time.Second, "Lock delay")
	cmd.Flags().String("session", "", "Previously created session to use for lock")

	addDatacenterOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func kvLock(cmd *cobra.Command, args []string) error {
	var lockOpts *consulapi.KVPair

	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	// Work around a Consul API bug that ignores LockDelay == 0
	lockDelay := viper.GetDuration("lock-delay")
	if lockDelay == 0 {
		lockDelay = time.Nanosecond
	}

	client, err := newKv()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()
	queryOpts := queryOptions()
	queryOpts.WaitTime = 15 * time.Second

	sessionClient, err := newSession()
	if err != nil {
		return err
	}

	if viper.GetString("session") == "" {
		// Create the Consul session
		se, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
			Name:      "Session for consul-cli",
			LockDelay: lockDelay,
			Behavior:  viper.GetString("behavior"),
			TTL:       viper.GetString("ttl"),
		}, writeOpts)
		if err != nil {
			return err
		}

		viper.Set("session", se)
		viper.Set("__clean_session", true)
	}

	session := viper.GetString("session")

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(viper.GetString("ttl"), session, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

WAIT:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		destroySession(sessionClient)
		return err
	}

	locked := false

	if kv != nil && kv.Session == session {
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
			Session: session,
		}
	} else {
		lockOpts = &consulapi.KVPair{
			Key:     kv.Key,
			Flags:   kv.Flags,
			Value:   kv.Value,
			Session: session,
		}
	}

	// Try to acquire the lock
	locked, _, err = client.Acquire(lockOpts, nil)
	if err != nil {
		destroySession(sessionClient)
		return err
	}

	if !locked {
		select {
		case <-time.After(5 * time.Second):
			goto WAIT
		}
	}

HELD:
	fmt.Println(session)

	return nil
}

// Destroy the session on error. Only performed when
// __clean_session == true
//
func destroySession(s *consulapi.Session) error {
	if viper.GetBool("__clean_session") {
		session := viper.GetString("session")
		writeOpts := writeOptions()
		_, err := s.Destroy(session, writeOpts)
		if err != nil {
			return fmt.Errorf("Session not destroyed: %s", session)
		}
	}

	return nil
}

