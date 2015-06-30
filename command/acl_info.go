package command

import (
	"encoding/json"
	"strings"
)

type ACLInfoCommand struct {
	Meta
}

func (c *ACLInfoCommand) Help() string {
	helpText := `
Usage: consul-cli acl-info [options] id

  Query information about an ACL token

Options:

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *ACLInfoCommand) Run(args []string) int {
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

	queryOpts := c.QueryOptions()
	acl, _, err := client.Info(id, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	jsonRaw, err := json.MarshalIndent(acl, "", "  ")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(string(jsonRaw))

	return 0
}

func (c *ACLInfoCommand) Synopsis() string {
	return "Query an ACL token"
}
