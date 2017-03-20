// +build consul8
//

package action

import (
	"flag"
)

type operatorAutopilotGet struct {
	stale bool

	*config
}

func OperatorAutopilotGetAction() Action {
	return &operatorAutopilotGet{
		config: &gConfig,
	}
}

func (o *operatorAutopilotGet) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&o.stale, "stale", false, "Read the raft configuration from any Consul server")

	o.addOutputFlags(f, false)
	o.addDatacenterFlag(f)

	return f
}

func (o *operatorAutopilotGet) Run(args []string) error {
	client, err := o.newOperator()
	if err != nil {
		return err
	}

	queryOpts := o.queryOptions()

	rc, err := client.AutopilotGetConfiguration(queryOpts)
	if err != nil {
		return err
	}

	return o.Output(rc)
}
