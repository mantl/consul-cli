package action

import (
	"flag"
)

type operatorKeyringList struct {
	*config
}

func OperatorKeyringListAction() Action {
	return &operatorKeyringList{
		config: &gConfig,
	}
}

func (o *operatorKeyringList) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (o *operatorKeyringList) Run(args []string) error {
	client, err := o.newOperator()
	if err != nil {
		return err
	}

	queryOpts := o.queryOptions()
	r, err := client.KeyringList(queryOpts)
	if err != nil {
		return err
	}

	return o.Output(r)
}
