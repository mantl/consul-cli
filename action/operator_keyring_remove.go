package action

import (
	"flag"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type operatorKeyringRemove struct {
	*config
}

func OperatorKeyringRemoveAction() Action {
	return &operatorKeyringRemove{
		config: &gConfig,
	}
}

func (o *operatorKeyringRemove) CommandFlags() *flag.FlagSet {
	return o.newFlagSet(FLAG_NONE)
}

func (o *operatorKeyringRemove) Run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("At least one gossip key must be specified")
	}

	client, err := o.newOperator()
	if err != nil {
		return err
	}

	writeOpts := o.writeOptions()

	var result error

	for _, k := range args {
		if err := client.KeyringRemove(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
