package action

import (
	"flag"
	"fmt"
)

type operatorKeyringUse struct {
	*config
}
	
func OperatorKeyringUseAction() Action {
	return &operatorKeyringUse{
		config: &gConfig,
	}
}

func (o *operatorKeyringUse) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (o *operatorKeyringUse) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Encryption key must be specified")
	}

	client, err := o.newOperator()
	if err != nil {
		return err
	}

	writeOpts := o.writeOptions()

	return client.KeyringUse(args[0], writeOpts)
}
