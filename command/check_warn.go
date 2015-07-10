package command

import (
	"strings"
)

type CheckWarnCommand struct {
	Meta
	note		string
}

func (c *CheckWarnCommand) Help() string {
	helpText := `
Usage: consul-cli check-warn [options] checkId

  Mark a local check as warning

Options:
` + c.ConsulHelp() + 
`  --note			Message to associate with check status
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *CheckWarnCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.StringVar(&c.note, "note", "", "")
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
	checkId := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Agent()
	err = client.WarnTTL(checkId, c.note)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *CheckWarnCommand) Synopsis() string {
	return "Mark a local check as warning"
}
