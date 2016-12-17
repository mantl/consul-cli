package commands

import (
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Delete functions

func newKvDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <path>",
		Short: "Delete a given path from the K/V",
		Long:  "Delete a given path from the K/V",
		RunE:  kvDelete,
	}

	cmd.Flags().String("modifyindex", "", "Perform a Check-and-Set delete")
	cmd.Flags().Bool("recurse", false, "Perform a recursive delete")
	addDatacenterOption(cmd)

	return cmd
}

func kvDelete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	switch {
	case viper.GetBool("recurse"):
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			return err
		}
	case viper.GetString("modifyindex") != "":
		m, err := strconv.ParseUint(viper.GetString("modifyindex"), 0, 64)
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

