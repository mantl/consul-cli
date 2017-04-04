package action

import (
	"flag"
	"fmt"
)

type aclInfo struct {
	*config
}

func AclInfoAction() Action {
	return &aclInfo{
		config: &gConfig,
	}
}

func (a *aclInfo) CommandFlags() *flag.FlagSet {
	f := newFlagSet()
	a.addOutputFlags(f, false)

	return f
}

func (a aclInfo) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	client, err := a.newACL()
	if err != nil {
		return err
	}

	queryOpts := a.queryOptions()
	acl, _, err := client.Info(id, queryOpts)
	if err != nil {
		return err
	}

	return a.Output(acl)
}
