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
	f := newFlagSet()

	k.addConsistencyFlags(f)
	k.addDatacenterFlag(f)
	k.addWaitIndexFlag(f)
	k.addOutputFlags(f, true)

	return f
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

