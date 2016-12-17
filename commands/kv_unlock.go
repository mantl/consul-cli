package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newKvUnlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock <path>",
		Short: "Release a lock on a given path",
		Long:  "Release a lock on a given path",
		RunE:  kvUnlock,
	}

	cmd.Flags().String("session", "", "Session ID of the lock holder. Required")
	cmd.Flags().Bool("no-destroy", false, "Do not destroy the session when complete")
	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)

	return cmd
}

func kvUnlock(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	session := viper.GetString("session")
	if session == "" {
		return fmt.Errorf("Session ID must be provided")
	}

	client, err := newKv()
	if err != nil {
		return err
	}

	sessionClient, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	kv, _, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		return fmt.Errorf("Node '%s' does not exist", path)
	}

	if kv.Session != session {
		return fmt.Errorf("Session not lock holder")
	}

	writeOpts := writeOptions()

	success, _, err := client.Release(kv, writeOpts)
	if err != nil {
		return err
	}

	if !viper.GetBool("no-destroy") {
		_, err = sessionClient.Destroy(session, writeOpts)
		if err != nil {
			return err
		}
	}

	if !success {
		return fmt.Errorf("Failed unlocking path")
	}

	return nil
}

