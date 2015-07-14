package command

import (
	"strings"
)

type AgentChecksCommand struct {
	Meta
}

func (c *AgentChecksCommand) Help() string {
	helpText := `
Usage: consul-cli agent-checks [options]

  Get the checks the agent is managing

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *AgentChecksCommand) Run(args []string) int {
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
	cs, err := client.Checks()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(cs, true)

	return 0
}

func (c *AgentChecksCommand) Synopsis() string {
	return "Get the checks the agent is managing"
}
