package command

import (
	"strings"
)

type SessionInfoCommand struct {
	Meta
}

func (c *SessionInfoCommand) Help() string {
	helpText := `
Usage: consul-cli session-info [options] sessionId

  Get information on a session

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *SessionInfoCommand) Run(args []string) int {
	c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Session ID must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	sessionid := extra[0]

	client, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	queryOpts := c.QueryOptions()
	sessionClient := client.Session()

	s, _, err := sessionClient.Info(sessionid, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(s, true)

	return 0
}

func (c *SessionInfoCommand) Synopsis() string {
	return "Get information on a session"
}
