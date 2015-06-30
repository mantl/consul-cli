package command

import (
	"encoding/json"
	"strings"
)

type ACLListCommand struct {
	Meta
}

func (c *ACLListCommand) Help() string {
	helpText := `
Usage: consul-cli list [options]

  List all active ACL tokens.

Options:

` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *ACLListCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	consul, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.ACL()

	queryOpts := c.QueryOptions()
	acls, _, err := client.List(queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	jsonRaw, err := json.MarshalIndent(acls, "", "  ")
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(string(jsonRaw))

	return 0
}

func (c *ACLListCommand) Synopsis() string {
	return "List a value"
}
