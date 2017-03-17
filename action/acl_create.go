package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type aclCreate struct {
	management bool
	name       string
	rules      []string
	*config
}

func AclCreateAction() Action {
	return &aclCreate{
		config: &gConfig,
	}
}

func (a *aclCreate) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&a.management, "management", false, "Create a management token")
	f.StringVar(&a.name, "name", "", "Name of the ACL")
	f.Var(newStringSliceValue(&a.rules), "rule", "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	a.addRawFlag(f)

	return f
}

func (a *aclCreate) Run(args []string) error {
	client, err := a.newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = a.name

	switch {
	case len(args) == 1:
		entry.ID = args[0]
	case len(args) > 1:
		return fmt.Errorf("Only one ACL identified can be specified")
	}

	if a.management {
		entry.Type = consulapi.ACLManagementType
	} else {
		entry.Type = consulapi.ACLClientType
	}

	if a.raw.isSet() {
		rules, err := a.raw.readString()
		if err != nil {
			return err
		}
		entry.Rules = rules
	} else {
		rules, err := getRulesString(a.rules)
		if err != nil {
			return err
		}
		entry.Rules = rules
	}

	writeOpts := a.writeOptions()
	id, _, err := client.Create(entry, writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}
