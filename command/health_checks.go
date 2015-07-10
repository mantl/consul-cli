package command

import (
	"strings"
)

type HealthChecksCommand struct {
	Meta
}

func (c *HealthChecksCommand) Help() string {
	helpText := `
Usage: consul-cli health-checks [options] service

  Get the health checks for a service

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *HealthChecksCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
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
	service := extra[0]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	queryOpts := c.QueryOptions()
	healthClient := client.Health()

	h, _, err := healthClient.Checks(service, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(h, true)

	return 0
}

func (c *HealthChecksCommand) Synopsis() string {
	return "Get the health checks for a service"
}
