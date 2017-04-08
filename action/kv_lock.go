package action

import (
	"flag"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type kvLock struct {
	behavior     string
	ttl          string
	lockDelay    time.Duration
	session      string
	cleanSession bool

	*config
}

func KvLockAction() Action {
	return &kvLock{
		config: &gConfig,
	}
}

func (k *kvLock) CommandFlags() *flag.FlagSet {
	f := k.newFlagSet(FLAG_DATACENTER, FLAG_CONSISTENCY)

	f.StringVar(&k.behavior, "behavior", "release", "Lock behavior. One of 'release' or 'delete'")
	f.StringVar(&k.ttl, "ttl", "", "Lock time to live")
	f.DurationVar(&k.lockDelay, "lock-delay", 15*time.Second, "Lock delay")
	f.StringVar(&k.session, "session", "", "Previously created session to use for lock")

	return f
}

func (k *kvLock) Run(args []string) error {
	var lockOpts *consulapi.KVPair

	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	// Work around a Consul API bug that ignores LockDelay == 0
	if k.lockDelay == 0 {
		k.lockDelay = time.Nanosecond
	}

	client, err := k.newKv()
	if err != nil {
		return err
	}

	writeOpts := k.writeOptions()
	queryOpts := k.queryOptions()
	queryOpts.WaitTime = 15 * time.Second

	sessionClient, err := k.newSession()
	if err != nil {
		return err
	}

	if k.session == "" {
		// Create the Consul session
		se, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
			Name:      "Session for consul-cli",
			LockDelay: k.lockDelay,
			Behavior:  k.behavior,
			TTL:       k.ttl,
		}, writeOpts)
		if err != nil {
			return err
		}

		k.session = se
		k.cleanSession = true
	}

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(k.ttl, k.session, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

WAIT:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		k.destroySession(sessionClient)
		return err
	}

	locked := false

	if kv != nil && kv.Session == k.session {
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
			Session: k.session,
		}
	} else {
		lockOpts = &consulapi.KVPair{
			Key:     kv.Key,
			Flags:   kv.Flags,
			Value:   kv.Value,
			Session: k.session,
		}
	}

	// Try to acquire the lock
	locked, _, err = client.Acquire(lockOpts, nil)
	if err != nil {
		k.destroySession(sessionClient)
		return err
	}

	if !locked {
		select {
		case <-time.After(5 * time.Second):
			goto WAIT
		}
	}

HELD:
	fmt.Println(k.session)

	return nil
}

// Destroy the session on error. Only performed when
// k.cleanSession == true
//
func (k *kvLock) destroySession(s *consulapi.Session) error {
	if k.cleanSession {
		writeOpts := k.writeOptions()
		_, err := s.Destroy(k.session, writeOpts)
		if err != nil {
			return fmt.Errorf("Session not destroyed: %s", k.session)
		}
	}

	return nil
}
