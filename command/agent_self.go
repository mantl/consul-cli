package command

import (
	"strings"
)

type AgentSelfCommand struct {
	Meta
}

func (c *AgentSelfCommand) Help() string {
	helpText := `
Usage: consul-cli agent-self [options]

  Get the node configuration

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *AgentSelfCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()
	config, err := client.Self()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *AgentSelfCommand) Synopsis() string {
	return "Get the node configuration"
}
