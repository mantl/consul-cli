package command

import (
	"strings"
)

type SessionRenewCommand struct {
	Meta
}

func (c *SessionRenewCommand) Help() string {
	helpText := `
Usage: consul-cli session-renew [options] sessionId

  Renew the given session

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *SessionRenewCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(true)
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

	writeOpts := c.WriteOptions()
	sessionClient := client.Session()

	s, _, err := sessionClient.Renew(sessionid, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if s != nil {
		c.OutputJSON(s, true)
	}


	return 0
}

func (c *SessionRenewCommand) Synopsis() string {
	return "Renew the given session"
}
