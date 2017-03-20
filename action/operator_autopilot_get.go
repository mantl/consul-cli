// +build consul8
//

package action

import (
	"flag"
)

type operatorAutopilotGet struct {
	*config
}

func OperatorAutopilotGetAction() Action {
	return &operatorAutopilotGet{
		config: &gConfig,
	}
}

func (o *operatorAutopilotGet) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	o.addOutputFlags(f, false)
	o.addDatacenterFlag(f)
	o.addStaleFlag(f)

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
