package command

import (
	"strings"
)

type AgentJoinCommand struct {
	Meta
	joinWAN			bool
}

func (c *AgentJoinCommand) Help() string {
	helpText := `
Usage: consul-cli agent-join [options] nodeName

  Trigger the agent to join a node

Options:
` + c.ConsulHelp() + 
`  --wan				Attempt to join WAN pool
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *AgentJoinCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.BoolVar(&c.joinWAN, "wan", false, "")
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
	nodeName := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()
	err = client.Join(nodeName, c.joinWAN)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *AgentJoinCommand) Synopsis() string {
	return "Trigger the local agent to join a node"
}
