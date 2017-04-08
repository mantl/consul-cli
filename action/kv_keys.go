package action

import (
	"flag"
	"fmt"
)

type kvKeys struct {
	separator string

	*config
}

func KvKeysAction() Action {
	return &kvKeys{
		config: &gConfig,
	}
}

func (k *kvKeys) CommandFlags() *flag.FlagSet {
	f := k.newFlagSet(FLAG_DATACENTER, FLAG_CONSISTENCY, FLAG_BLOCKING)

	f.StringVar(&k.separator, "separator", "", "List keys only up to a given separator")

	return f
}

func (k *kvKeys) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	client, err := k.newKv()
	if err != nil {
		return err
	}

	queryOpts := k.queryOptions()
	data, _, err := client.Keys(path, k.separator, queryOpts)
	if err != nil {
		return err
	}

	k.output.template = kv_outputTemplate

	return k.Output(data)
}

var kv_outputTemplate = `{{range .}}{{.}}
{{end}}`
