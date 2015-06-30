package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type ACLUpdateCommand struct {
	Meta
	ConfigRules	[]*ConfigRule
}

func (c *ACLUpdateCommand) Help() string {
	helpText := `
Usage: consul-cli acl-update [options] id

  Update an ACL. Will be created if it doesn't exist.

Options:

` + c.ConsulHelp() +
`  --management			Update type to 'management'
				(default: false)
  --name			Name of the ACL
				(default: not set)
  --rule='type:path:policy'	Rule to create. Can be multiple rules on a command line
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *ACLUpdateCommand) Run(args []string) int {
	var isManagement bool
	var aclName string

	flags := c.Meta.FlagSet()
	flags.StringVar(&aclName, "name", "", "")
	flags.BoolVar(&isManagement, "management", false, "")

	flags.Var((funcVar)(func(s string) error {
		t, err := ParseRuleConfig(s)
		if err != nil {
			return err
		}

		if c.ConfigRules == nil {
			c.ConfigRules = make([]*ConfigRule, 0, 1)
		}

		c.ConfigRules = append(c.ConfigRules, t)
		return nil
	}), "rule", "")
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

	var entry *consulapi.ACLEntry

	if isManagement {
		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	aclName,
			Type:	consulapi.ACLManagementType,
		}
	} else {
		rules, err := GetRulesString(c.ConfigRules)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	aclName,
			Type:	consulapi.ACLClientType,
			Rules:	rules,
		}

	}

	writeOpts := c.WriteOptions()
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *ACLUpdateCommand) Synopsis() string {
	return "Update an ACL"
}
