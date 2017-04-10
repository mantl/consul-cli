package action

import (
	"flag"
	"fmt"
)

type aclDestroy struct {
	*config
}

func AclDestroyAction() Action {
	return &aclDestroy{
		config: &gConfig,
	}
}

func (a *aclDestroy) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_NONE)
}

func (a *aclDestroy) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}
	id := args[0]

	client, err := a.newACL()
	if err != nil {
		return err
	}

	writeOpts := a.writeOptions()
	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
