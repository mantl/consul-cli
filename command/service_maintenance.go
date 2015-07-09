package command

import (
	"strings"
)

type ServiceMaintenanceCommand struct {
	Meta
	enabled			bool
	reason			string
}

func (c *ServiceMaintenanceCommand) Help() string {
	helpText := `
Usage: consul-cli service-maintenance [options] serviceId

  Manage maintenance mode of a service

Options:
` + c.ConsulHelp() + 
`  --enabled			Boolean value for maintenance mode
				(default: true)
  --reason			Reason for entering maintenance mode
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *ServiceMaintenanceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.BoolVar(&c.enabled, "enabled", true, "")
	flags.StringVar(&c.reason, "reason", "", "")
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
	serviceId := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()

	if c.enabled {
		err = client.EnableServiceMaintenance(serviceId, c.reason)
	} else {
		err = client.DisableServiceMaintenance(serviceId)
	}
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *ServiceMaintenanceCommand) Synopsis() string {
	return "Manage maintenance mode on a service"
}
