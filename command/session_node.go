package command

import (
	"strings"
)

type SessionNodeCommand struct {
	Meta
}

func (c *SessionNodeCommand) Help() string {
	helpText := `
Usage: consul-cli session-node [options] nodeName

  Get active sessions for a node

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *SessionNodeCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(true)
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Node name must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	node := extra[0]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	queryOpts := c.QueryOptions()
	sessionClient := client.Session()

	s, _, err := sessionClient.Node(node, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(s, true)

	return 0
}

func (c *SessionNodeCommand) Synopsis() string {
	return "Get active sessions for a node"
}
