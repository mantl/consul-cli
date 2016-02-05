package commands

import (
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type KvDeleteOptions struct {
	ModifyIndex string
	DoRecurse   bool
}

func (k *Kv) AddDeleteSub(cmd *cobra.Command) {
	kdo := &KvDeleteOptions{}

	deleteCmd := &cobra.Command{
		Use:   "delete <path>",
		Short: "Delete a given path from the K/V",
		Long:  "Delete a given path from the K/V",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Delete(args, kdo)
		},
	}

	oldDeleteCmd := &cobra.Command{
		Use:        "kv-delete <path>",
		Short:      "Delete a given path from the K/V",
		Long:       "Delete a given path from the K/V",
		Deprecated: "Use kv delete",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Delete(args, kdo)
		},
	}

	deleteCmd.Flags().StringVar(&kdo.ModifyIndex, "modifyindex", "", "Perform a Check-and-Set delete")
	deleteCmd.Flags().BoolVar(&kdo.DoRecurse, "recurse", false, "Perform a recursive delete")
	k.AddDatacenterOption(deleteCmd)

	oldDeleteCmd.Flags().StringVar(&kdo.ModifyIndex, "modifyindex", "", "Perform a Check-and-Set delete")
	oldDeleteCmd.Flags().BoolVar(&kdo.DoRecurse, "recurse", false, "Perform a recursive delete")
	k.AddDatacenterOption(oldDeleteCmd)

	cmd.AddCommand(deleteCmd)

	k.AddCommand(oldDeleteCmd)
}

func (k *Kv) Delete(args []string, kdo *KvDeleteOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Key path must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one key path allowed")
	}
	path := args[0]

	client, err := k.KV()
	if err != nil {
		return err
	}

	writeOpts := k.WriteOptions()

	switch {
	case kdo.DoRecurse:
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			return err
		}
	case kdo.ModifyIndex != "":
		m, err := strconv.ParseUint(kdo.ModifyIndex, 0, 64)
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
