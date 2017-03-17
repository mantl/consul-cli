package action

import (
	"flag"
	"fmt"
)

type aclClone struct {
	*config
}

func AclCloneAction() Action {
	return &aclClone{
		config: &gConfig,
	}
}

func (a *aclClone) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (a *aclClone) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}

	client, err := a.consul.newACL()
	if err != nil {
		return err
	}

	writeOpts := a.consul.writeOptions()

	newid, _, err := client.Clone(args[0], writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(newid)

	return nil
}
