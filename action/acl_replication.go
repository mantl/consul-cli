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
	return a.newFlagSet(FLAG_OUTPUT, FLAG_DATACENTER)
}

func (a *aclReplication) Run(args []string) error {
	return fmt.Errorf("ACL replication status not available in Consul API")
}
