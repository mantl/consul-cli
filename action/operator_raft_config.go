package action

import (
	"flag"
)

type operatorRaftConfig struct {
	stale bool

	*config
}

func OperatorRaftConfigAction() Action {
	return &operatorRaftConfig{
		config: &gConfig,
	}
}

func (o *operatorRaftConfig) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&o.stale, "stale", false, "Read the raft configuration from any Consul server")

	o.addOutputFlags(f, false)
	o.addDatacenterFlag(f)

	return f
}

func (o *operatorRaftConfig) Run(args []string) error {
	client, err := o.newOperator()
	if err != nil {
		return err
	}

	queryOpts := o.queryOptions()

	rc, err := client.RaftGetConfiguration(queryOpts)
	if err != nil {
		return err
	}

	return o.Output(rc)
}
