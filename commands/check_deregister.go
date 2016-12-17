package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Deregister functions

func newCheckDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister",
		Short: "Remove a check from the agent",
		Long:  "Remove a check from the agent",
		RunE:  checkDeregister,
	}

	return cmd
}

func checkDeregister(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.CheckDeregister(checkId)
	if err != nil {
		return err
	}

	return nil
}
