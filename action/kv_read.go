package action

import (
	"flag"
	"fmt"
)

type kvRead struct {
	recurse bool

	*config
}

func KvReadAction() Action {
	return &kvRead{
		config: &gConfig,
	}
}

func (k *kvRead) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&k.recurse, "recurse", false, "Perform a recursive read")

	k.addConsistencyFlags(f)
	k.addDatacenterFlag(f)
	k.addOutputFlags(f, true)

	return f
}

func (k *kvRead) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	client, err := k.newKv()
	if err != nil {
		return err
	}

	queryOpts := k.queryOptions()

	if k.recurse {
		kvlist, _, err := client.List(path, queryOpts)
		if err != nil {
			return err
		}

		if kvlist == nil {
			return nil
		}

		return k.OutputKv(kvlist)
	} else {
		kv, _, err := client.Get(path, queryOpts)
		if err != nil {
			return err
		}

		if kv == nil {
			return nil
		}

		return k.OutputKv(kv)
	}
}

