package action

import (
	"flag"
	"fmt"
)

type operatorRaftDelete struct {
	*config
}

func OperatorRaftDeleteAction() Action {
	return &operatorRaftDelete{
		config: &gConfig,
	}
}

func (o *operatorRaftDelete) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	o.addDatacenterFlag(f)

	return f
}

func (o *operatorRaftDelete) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single address argument must be specified")
	}
	address := args[0]

	client, err := o.newOperator()
	if err != nil {
		return err
	}

	writeOpts := o.writeOptions()

	return client.RaftRemovePeerByAddress(address, writeOpts)
}
