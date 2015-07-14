package command

import (
	"strings"
)

type ACLDestroyCommand struct {
	Meta
}

func (c *ACLDestroyCommand) Help() string {
	helpText := `
Usage: consul-cli acl-destroy [options] token

  Destroy an ACL

Options:

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *ACLDestroyCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("ACL id must be specified")
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
	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *ACLDestroyCommand) Synopsis() string {
	return "Destroy an ACL"
}
