package action

import (
	"flag"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type operatorKeyringInstall struct {
	*config
}

func OperatorKeyringInstallAction() Action {
	return &operatorKeyringInstall{
		config: &gConfig,
	}
}

func (o * operatorKeyringInstall) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (o *operatorKeyringInstall) Run(args []string) error {
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
		if err := client.KeyringInstall(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
