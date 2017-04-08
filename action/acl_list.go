package action

import (
	"flag"
)

// List functions

type aclList struct {
	*config
}

func AclListAction() Action {
	return &aclList{
		config: &gConfig,
	}
}

func (a *aclList) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_OUTPUT, FLAG_CONSISTENCY, FLAG_BLOCKING)
}

func (a *aclList) Run(args []string) error {
	client, err := a.newACL()
	if err != nil {
		return err
	}

	queryOpts := a.queryOptions()
	acls, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return a.Output(acls)
}
