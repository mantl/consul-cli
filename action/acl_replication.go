package action

import (
	"flag"
	"fmt"
)

type aclReplication struct {
	*config
}

func AclReplicationAction() Action {
	return &aclReplication{
		config: &gConfig,
	}
}

func (a *aclReplication) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	a.addDatacenterFlag(f)
	a.addOutputFlags(f, false)

	return f
}

func (a *aclReplication) Run(args []string) error {
	return fmt.Errorf("ACL replication status not available in Consul API")
}
