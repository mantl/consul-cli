package command

import (
	"strings"
)

type ServiceDeregisterCommand struct {
	Meta
}

func (c *ServiceDeregisterCommand) Help() string {
	helpText := `
Usage: consul-cli service-deregister [options] serviceId

  Remove a service from the agent

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *ServiceDeregisterCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
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
	err = client.ServiceDeregister(serviceId)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *ServiceDeregisterCommand) Synopsis() string {
	return "Remove a service from the agent"
}
