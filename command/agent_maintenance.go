package command

import (
	"strings"
)

type AgentMaintenanceCommand struct {
	Meta
	enabled			bool
	reason			string
}

func (c *AgentMaintenanceCommand) Help() string {
	helpText := `
Usage: consul-cli agent-maintenance [options] 

  Manage node maintenance mode

Options:
` + c.ConsulHelp() + 
`  --enabled			Boolean value for maintenance mode
				(default: true)
  --reason			Reason for entering maintenance mode
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *AgentMaintenanceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.BoolVar(&c.enabled, "enabled", true, "")
	flags.StringVar(&c.reason, "reason", "", "")
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

	if c.enabled {
		err = client.EnableNodeMaintenance(c.reason)
	} else {
		err = client.DisableNodeMaintenance()
	}
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *AgentMaintenanceCommand) Synopsis() string {
	return "Manage node maintenance mode"
}
