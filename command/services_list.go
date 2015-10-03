package command

import (
	"strings"
)

type ServicesListCommand struct {
	Meta
}

func (c *ServicesListCommand) Help() string {
	helpText := `
Usage: consul-cli services-list

  List local services
`

	return strings.TrimSpace(helpText)
}

func (c *ServicesListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	consul, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()
	out, err := client.Services()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(out, true)

	return 0
}

func (c *ServicesListCommand) Synopsis() string {
	return "List local service"
}
