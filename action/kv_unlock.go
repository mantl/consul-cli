package action

import (
	"flag"
	"fmt"
)

type kvUnlock struct {
	session string
	noDestroy bool

	*config
}

func KvUnlockAction() Action {
	return &kvUnlock{
		config: &gConfig,
	}
}

func (k *kvUnlock) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&k.session, "session", "", "Session ID of the lock holder. Required")
	f.BoolVar(&k.noDestroy, "no-destroy", false, "Do not destroy the session when complete")

	k.addConsistencyFlags(f)
	k.addDatacenterFlag(f)

	return f
}

func (k *kvUnlock) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	if k.session == "" {
		return fmt.Errorf("Session ID must be provided")
	}

	client, err := k.newKv()
	if err != nil {
		return err
	}

	sessionClient, err := k.newSession()
	if err != nil {
		return err
	}

	queryOpts := k.queryOptions()

	kv, _, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		return fmt.Errorf("Node '%s' does not exist", path)
	}

	if kv.Session != k.session {
		return fmt.Errorf("Session not lock holder")
	}

	writeOpts := k.writeOptions()

	success, _, err := client.Release(kv, writeOpts)
	if err != nil {
		return err
	}

	if !k.noDestroy {
		_, err = sessionClient.Destroy(k.session, writeOpts)
		if err != nil {
			return err
		}
	}

	if !success {
		return fmt.Errorf("Failed unlocking path")
	}

	return nil
}

