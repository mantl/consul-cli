package action

import (
	"flag"
	"fmt"
)

type kvWatch struct {
	*config
}

func KvWatchAction() Action {
	return &kvWatch{
		config: &gConfig,
	}
}

func (k *kvWatch) CommandFlags() *flag.FlagSet {
	return k.newFlagSet(FLAG_DATACENTER, FLAG_KVOUTPUT, FLAG_CONSISTENCY, FLAG_BLOCKING)
}

func (k *kvWatch) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	client, err := k.newKv()
	if err != nil {
		return err
	}

	queryOpts := k.queryOptions()

RETRY:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		queryOpts.WaitIndex = meta.LastIndex
		goto RETRY
	}

	return k.OutputKv(kv)
}
