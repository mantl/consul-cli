package command

import (
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type SessionCreateCommand struct {
	Meta
	lockDelay		time.Duration
	nodeName		string
	name			string
	checks			[]string
	behavior		string
	ttl			time.Duration
}

func (c *SessionCreateCommand) Help() string {
	helpText := `
Usage: consul-cli session-list [options]

  Create a new session

Options: 

` + c.ConsulHelp() +
`  --lock-delay			Lock delay as a duration string
				(default: not set)
  --name			Session name
				(default: not set)
  --node			Node to register session
				(default: agent node name)
  --checks			Check to associate with session. Can be mulitple
				(default: not set)
  --behavior			Lock behavior when session is invalidated. One of
				release or delete
				(default: release)
  --ttl				Session Time To Live as a duration string
				(default: 15s)
`

	return strings.TrimSpace(helpText)
}

func (c *SessionCreateCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(true)
	flags.DurationVar(&c.lockDelay, "lock-delay", 0, "")
	flags.StringVar(&c.nodeName, "node", "", "")
	flags.StringVar(&c.name, "name", "", "")
	flags.DurationVar(&c.ttl, "ttl", 15 * time.Second, "")
	flags.Var((funcVar)(func(s string) error {
		if c.checks == nil {
			c.checks = make([]string, 0, 1)
		}

		c.checks = append(c.checks, s)
		return nil
	}), "checks", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	// Work around Consul API bug that drops LockDelay == 0
	if c.lockDelay == 0 {
		c.lockDelay = time.Nanosecond
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	writeOpts := c.WriteOptions()
	sessionClient := client.Session()

	s := consulapi.SessionEntry{
		Name:		c.name,
		Node:		c.nodeName,
		Checks:		c.checks,
		LockDelay:	c.lockDelay,
		Behavior:	c.behavior,
		TTL:		c.ttl.String(),
	}

	se, _, err := sessionClient.Create(&s, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(se)

	return 0
}

func (c *SessionCreateCommand) Synopsis() string {
	return "Create a new session"
}
