package command

import (
	"strings"
)

type AgentServicesCommand struct {
	Meta
}

func (c *AgentServicesCommand) Help() string {
	helpText := `
Usage: consul-cli agent-services [options]

  Get the services the agent is managing

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *AgentServicesCommand) Run(args []string) int {
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
	config, err := client.Services()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *AgentServicesCommand) Synopsis() string {
	return "Get the services the agent is managing"
}
