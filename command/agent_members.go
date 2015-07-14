package command

import (
	"strings"
)

type AgentMembersCommand struct {
	Meta
	wanFlag		bool
}

func (c *AgentMembersCommand) Help() string {
	helpText := `
Usage: consul-cli agent-members [options]

  Get the members as seen by the serf agent

Options:
` + c.ConsulHelp() + 
`  --wan				Get list of WAN members instead of LAN
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *AgentMembersCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
	flags.BoolVar(&c.wanFlag, "wan", false, "")
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
	ms, err := client.Members(c.wanFlag)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(ms, true)

	return 0
}

func (c *AgentMembersCommand) Synopsis() string {
	return "Get the members as seen by the serf agent"
}
