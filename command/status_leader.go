package command

import (
	"strings"
)

type StatusLeaderCommand struct {
	Meta
}

func (c *StatusLeaderCommand) Help() string {
	helpText := `
Usage: consul-cli status-leader [options]

  Get the current Raft leader

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *StatusLeaderCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	statusClient := client.Status()

	s, err := statusClient.Leader()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(s)

	return 0
}

func (c *StatusLeaderCommand) Synopsis() string {
	return "Get the current Raft leader"
}
