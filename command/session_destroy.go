package command

import (
	"strings"
)

type SessionDestroyCommand struct {
	Meta
}

func (c *SessionDestroyCommand) Help() string {
	helpText := `
Usage: consul-cli session-destroy [options] sessionId

  Destroy a session

Options: 

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *SessionDestroyCommand) Run(args []string) int {
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

	writeOpts := c.WriteOptions()
	sessionClient := client.Session()

	_, err = sessionClient.Destroy(sessionid, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *SessionDestroyCommand) Synopsis() string {
	return "Destroy a session"
}
