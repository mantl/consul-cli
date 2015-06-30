package command

import (
	"encoding/json"
	"strings"
)

type SessionListCommand struct {
	Meta
}

func (c *SessionListCommand) Help() string {
	helpText := `
Usage: consul-cli session-list [options]

  List active sessions for a datacenter

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *SessionListCommand) Run(args []string) int {
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

	queryOpts := c.QueryOptions()
	sessionClient := client.Session()

	sessions, _, err := sessionClient.List(queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	jsonRaw, err := json.MarshalIndent(sessions, "", "  ")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(string(jsonRaw))

	return 0
}

func (c *SessionListCommand) Synopsis() string {
	return "List active sessions for a datacenter"
}
