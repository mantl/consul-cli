package command

import (
	"strings"
)

type ACLCloneCommand struct {
	Meta
}

func (c *ACLCloneCommand) Help() string {
	helpText := `
Usage: consul-cli acl-clone [options] token

  Create a new token from an existing one

Options:

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *ACLCloneCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("ACL id must be provided")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	id := extra[0]

	consul, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.ACL()

	writeOpts := c.WriteOptions()
	newid, _, err := client.Clone(id, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(newid)

	return 0
}

func (c *ACLCloneCommand) Synopsis() string {
	return "Create a new token from an existing one"
}
