package command

import (
	"strings"
)

type HealthStateCommand struct {
	Meta
}

func (c *HealthStateCommand) Help() string {
	helpText := `
Usage: consul-cli health-checks [options] (any|unknown|passing|warning|critical)

  Get the checks in a given state

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *HealthStateCommand) Run(args []string) int {
	c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Check state must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	state := strings.ToLower(extra[0])

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	queryOpts := c.QueryOptions()
	healthClient := client.Health()

	h, _, err := healthClient.State(state, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(h, true)

	return 0
}

func (c *HealthStateCommand) Synopsis() string {
	return "Get the checks in a given state"
}
