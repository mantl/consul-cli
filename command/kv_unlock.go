package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

const (
	LockFlagValue = 0x2ddccbc058a50c18
)

type KVUnlockCommand struct {
	Meta
}

func (c *KVUnlockCommand) Help() string {
	helpText := `
Usage: consul-cli kv-unlock [options] path

  Release a lock on a given path

Options:

` + c.ConsulHelp() +
`  --session=<sessionId>		Session ID of the lock holder. Required
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *KVUnlockCommand) Run(args []string) int {
	var sessionId string

	flags := c.Meta.FlagSet()
	flags.StringVar(&sessionId, "session", "", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	if sessionId == "" {
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
	client := consul.KV()
	sessionClient := consul.Session()

	kv := new(consulapi.KVPair)
	kv.Key = path
	kv.Session = sessionId
	kv.Flags = LockFlagValue

	writeOpts := c.WriteOptions()

	success, _, err := client.Release(kv, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	_, err = sessionClient.Destroy(sessionId, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if !success {
		return 1
	}

	return 0
}

func (c *KVUnlockCommand) Synopsis() string {
	return "Unlock a node"
}
