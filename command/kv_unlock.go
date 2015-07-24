package command

import (
	"fmt"
	"strings"
)

type KVUnlockCommand struct {
	Meta
	session		string
	noDestroy	bool
}

func (c *KVUnlockCommand) Help() string {
	helpText := `
Usage: consul-cli kv-unlock [options] path

  Release a lock on a given path

Options:

` + c.ConsulHelp() +
`  --session=<sessionId>		Session ID of the lock holder. Required
				(default: not set)
  --no-destroy			Do not destroy the session when complete
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *KVUnlockCommand) Run(args []string) int {
	c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.StringVar(&c.session, "session", "", "")
	flags.BoolVar(&c.noDestroy, "no-destroy", false, "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if c.session == "" {
		c.UI.Error("Session ID must be provided")
		c.UI.Error("")
		c.UI.Error(c.Help())
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

	consul, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	kvClient := consul.KV()
	sessionClient := consul.Session()

	queryOpts := c.QueryOptions()

	kv, _, err := kvClient.Get(path, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if kv == nil {
		c.UI.Error(fmt.Sprintf("Node '%s' does not exist", path))
		return 1
	}

	if kv.Session != c.session {
		c.UI.Error("Session not lock holder")
		return 1
	}

	// clear the session
	kv.Session = ""

	writeOpts := c.WriteOptions()

	success, _, err := kvClient.Release(kv, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if !c.noDestroy {
		_, err = sessionClient.Destroy(c.session, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	}

	if !success {
		return 1
	}

	return 0
}

func (c *KVUnlockCommand) Synopsis() string {
	return "Unlock a node"
}
