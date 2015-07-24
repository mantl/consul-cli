package command

import (
	"fmt"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type KVLockCommand struct {
	Meta
	behavior	string
	ttl		string
	lockDelay	time.Duration
	session		string
	cleanSession	bool
}

func (c *KVLockCommand) Help() string {
	helpText := `
Usage: consul-cli kv-lock [options] path

  Acquire a lock on a given path

Options:

` + c.ConsulHelp() +
`  --behavior=release		Lock behavior. One of 'release' or 'delete'
				(default: release)
  --ttl				Lock time to live
				(default: not set)
  --lock-delay			Lock delay
				(default: 15s)
  --session			Previously created session to use for lock
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *KVLockCommand) Run(args []string) int {
	var lockOpts *consulapi.KVPair

	c.cleanSession = false

	c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.StringVar(&c.behavior, "behavior", "release", "")
	flags.StringVar(&c.ttl, "ttl", "", "")
	flags.DurationVar(&c.lockDelay, "lock-delay", time.Second * 15, "")
	flags.StringVar(&c.session, "session", "", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Work around a Consul API bug that ignores LockDelay == 0
	if c.lockDelay == 0 {
		c.lockDelay = time.Nanosecond
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Key path must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	path := extra[0]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	writeOpts := c.WriteOptions()
	queryOpts := c.QueryOptions()
	queryOpts.WaitTime = 15 * time.Second
	sessionClient := client.Session()
	kvClient := client.KV()

	if c.session == "" {
		// Create the Consul session
		se, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
						Name:		"Session for consul-cli",
						LockDelay:	c.lockDelay,
						Behavior:	c.behavior,
						TTL:		c.ttl,
						}, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		c.session = se
		c.cleanSession = true
	}

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(c.ttl, c.session, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

WAIT:
	kv, meta, err := kvClient.Get(path, queryOpts)
	if err != nil {
		c.destroySession(sessionClient)
		c.UI.Error(err.Error())
		return 1
	}

	locked := false

	if kv != nil && kv.Session == c.session {
		goto HELD
	}
	if kv != nil && kv.Session != "" {
		queryOpts.WaitIndex = meta.LastIndex
		goto WAIT
	}

	// Node doesn't already exist
	if kv == nil {
		lockOpts = &consulapi.KVPair {
			Key:		path,
			Session:	c.session,
		}
	} else {
		lockOpts = &consulapi.KVPair{
			Key:		kv.Key,
			Flags:		kv.Flags,
			Value:		kv.Value,
			Session:	c.session,
			}
	}

	// Try to acquire the lock
	locked, _, err = kvClient.Acquire(lockOpts, nil)
	if err != nil {
		c.destroySession(sessionClient)
		c.UI.Error(err.Error())
		return 1
	}

	if !locked {
		select {
		case <-time.After(5 * time.Second):
			goto WAIT
		}
	}

HELD:
	c.UI.Output(c.session)

	return 0
}

// Destroy the session on error. Only performed when
// c.cleanSession == true
//
func (c *KVLockCommand) destroySession(s *consulapi.Session) {
	if c.cleanSession {
		writeOpts := c.WriteOptions()
		_, err := s.Destroy(c.session, writeOpts)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Session not destroyed: %s", c.session))
		}
	}
}

func (c *KVLockCommand) Synopsis() string {
	return "Lock a node"
}
