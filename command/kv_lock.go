package command

import (
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type KVLockCommand struct {
	Meta
}

func (c *KVLockCommand) Help() string {
	helpText := `
Usage: consul-cli kv-lock [options] path

  Acquire a lock on a given path

Options:

` + c.ConsulHelp() +
`  --behavior=release		KVLock behavior. One of 'release' or 'delete'
				(default: release)
  --ttl=15s			KVLock time to live
				(default: 15s)
  --lock-delay=5s		KVLock delay
				(default: 5s)
`

	return strings.TrimSpace(helpText)
}

func (c *KVLockCommand) Run(args []string) int {
	var behavior string
	var ttl string
	var lockDelay time.Duration

	flags := c.Meta.FlagSet()
	flags.StringVar(&behavior, "behavior", "release", "")
	flags.StringVar(&ttl, "ttl", "15s", "")
	flags.DurationVar(&lockDelay, "lock-delay", 5 * time.Second, "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
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
	sessionClient := client.Session()

	// Create the Consul session
	sessionId, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
						Name:		"Session for consul-cli",
						LockDelay:	lockDelay,
						Behavior:	behavior,
						TTL:		ttl,
						}, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(ttl, sessionId, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

	// Create the KVLock Structure
	lockOpts := consulapi.LockOptions{
		Key:		path,
		Session:	sessionId,
		}
	l, err := client.LockOpts(&lockOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	_, err = l.Lock(nil)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(sessionId)

	return 0
}

func (c *KVLockCommand) Synopsis() string {
	return "Lock a node"
}
