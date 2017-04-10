package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type aclUpdate struct {
	management bool
	name       string
	rules      []string
	*config
}

func AclUpdateAction() Action {
	return &aclUpdate{
		config: &gConfig,
	}
}

func (a *aclUpdate) CommandFlags() *flag.FlagSet {
	f := a.newFlagSet(FLAG_RAW)

	f.BoolVar(&a.management, "management", false, "Create a management token")
	f.StringVar(&a.name, "name", "", "Name of the ACL")
	f.Var(newStringSliceValue(&a.rules), "rule", "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	return f
}

func (a *aclUpdate) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	client, err := a.newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = a.name
	entry.ID = id

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
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
