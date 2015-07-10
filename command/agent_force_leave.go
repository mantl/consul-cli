package command

import (
	"strings"
)

type AgentForceLeaveCommand struct {
	Meta
}

func (c *AgentForceLeaveCommand) Help() string {
	helpText := `
Usage: consul-cli agent-force-leave [options] nodeName

  Force the removal of a node

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *AgentForceLeaveCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Service name must be specified")
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
	err = client.ForceLeave(nodeName)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *AgentForceLeaveCommand) Synopsis() string {
	return "Force the removal of a node"
}
