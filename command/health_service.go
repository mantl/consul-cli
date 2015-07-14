package command

import (
	"strings"
)

type HealthServiceCommand struct {
	Meta
	tag		string
	passingOnly	bool
}

func (c *HealthServiceCommand) Help() string {
	helpText := `
Usage: consul-cli health-service [options] serviceName

  Get the nodes and health info for a service

Options: 

` + c.ConsulHelp() + 
`  --tag				Service tag to filter on
				(default: not set)
  --passing			Only return passing checks
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *HealthServiceCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(true)
	flags.StringVar(&c.tag, "tag", "", "")
	flags.BoolVar(&c.passingOnly, "passing", false, "")
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
	service := extra[0]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	queryOpts := c.QueryOptions()
	healthClient := client.Health()

	h, _, err := healthClient.Service(service, c.tag, c.passingOnly, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(h, true)

	return 0
}

func (c *HealthServiceCommand) Synopsis() string {
	return "Get nodes and health info for a service"
}
