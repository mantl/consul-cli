package command

import (
	"strings"
)

type StatusPeersCommand struct {
	Meta
}

func (c *StatusPeersCommand) Help() string {
	helpText := `
Usage: consul-cli status-peers [options]

  Get the current Raft peer-set

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *StatusPeersCommand) Run(args []string) int {
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

	s, err := statusClient.Peers()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	for _, i := range s {
		c.UI.Output(i)
	}

	return 0
}

func (c *StatusPeersCommand) Synopsis() string {
	return "Get the current Raft peer set"
}
