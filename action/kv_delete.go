package action

import (
	"flag"
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type kvDelete struct {
	modifyIndex string
	recurse     bool

	*config
}

func KvDeleteAction() Action {
	return &kvDelete{
		config: &gConfig,
	}
}

func (k *kvDelete) CommandFlags() *flag.FlagSet {
	f := k.newFlagSet(FLAG_DATACENTER)

	f.StringVar(&k.modifyIndex, "modifyindex", "", "Perform a Check-and-Set delete")
	f.BoolVar(&k.recurse, "recurse", false, "Perform a recursive delete")

	return f
}

func (k *kvDelete) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	client, err := k.newKv()
	if err != nil {
		return err
	}

	writeOpts := k.writeOptions()

	switch {
	case k.recurse:
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			return err
		}
	case k.modifyIndex != "":
		m, err := strconv.ParseUint(k.modifyIndex, 0, 64)
		if err != nil {
			return err
		}
		kv := consulapi.KVPair{
			Key:         path,
			ModifyIndex: m,
		}

		success, _, err := client.DeleteCAS(&kv, writeOpts)
		if err != nil {
			return err
		}

		if !success {
			return fmt.Errorf("Failed deleting")
		}
	default:
		_, err := client.Delete(path, writeOpts)
		if err != nil {
			return err
		}
	}

	return nil
}
