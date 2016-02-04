package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KvUnlockOptions struct {
	Session		string
	NoDestroy	bool
}

func (k *Kv) AddUnlockSub(cmd *cobra.Command) {
	kuo := &KvUnlockOptions{}

	unlockCmd := &cobra.Command{
		Use: "unlock <path>",
		Short: "Release a lock on a given path",
		Long: "Release a lock on a given path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Unlock(args, kuo)
		},
	}

	oldUnlockCmd := &cobra.Command{
		Use: "kv-unlock <path>",
		Short: "Release a lock on a given path",
		Long: "Release a lock on a given path",
		Deprecated: "Use kv unlock",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Unlock(args, kuo)
		},
	}

	unlockCmd.Flags().StringVar(&kuo.Session, "session", "", "Session ID of the lock holder. Required")
	unlockCmd.Flags().BoolVar(&kuo.NoDestroy, "no-destroy", false, "Do not destroy the session when complete")
	k.AddDatacenterOption(unlockCmd)

	oldUnlockCmd.Flags().StringVar(&kuo.Session, "session", "", "Session ID of the lock holder. Required")
	oldUnlockCmd.Flags().BoolVar(&kuo.NoDestroy, "no-destroy", false, "Do not destroy the session when complete")
	k.AddDatacenterOption(oldUnlockCmd)

	cmd.AddCommand(unlockCmd)

	k.AddCommand(oldUnlockCmd)
}

func (k *Kv) Unlock(args []string, kuo *KvUnlockOptions) error {
	if kuo.Session == "" {
		return fmt.Errorf("Session ID must be provided")
	}

	switch {
	case len(args) == 0:
		return fmt.Errorf("Key path must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one key path allowed")
	}
	path := args[0]

	consul, err := k.Client()
	if err != nil {
		return err
	}
	kvClient := consul.KV()
	sessionClient := consul.Session()

	queryOpts := k.QueryOptions()

	kv, _, err := kvClient.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		return fmt.Errorf("Node '%s' does not exist", path)
	}

	if kv.Session != kuo.Session {
		return fmt.Errorf("Session not lock holder")
	}

	writeOpts := k.WriteOptions()

	success, _, err := kvClient.Release(kv, writeOpts)
	if err != nil {
		return err
	}

	if !kuo.NoDestroy {
		_, err = sessionClient.Destroy(kuo.Session, writeOpts)
		if err != nil {
			return err
		}
	}

	if !success {
		return fmt.Errorf("Failed unlocking path")
	}

	return nil
}

