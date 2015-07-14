package command

import (
	"strings"
)

type CheckDeregisterCommand struct {
	Meta
}

func (c *CheckDeregisterCommand) Help() string {
	helpText := `
Usage: consul-cli check-deregister [options] checkId

  Remove a check from the agent

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *CheckDeregisterCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Check name must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	checkId := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()
	err = client.CheckDeregister(checkId)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *CheckDeregisterCommand) Synopsis() string {
	return "Remove a check from the agent"
}
