package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Raft configuration functions

func newOperatorRaftConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect the Raft configuration",
		Long:  "Inspect the Raft configuration",
		RunE:  operatorRaftConfig,
	}

	addTemplateOption(cmd)
	addDatacenterOption(cmd)

	cmd.Flags().Bool("stale", false, "Read the raft configuration from any Consul server")

	return cmd
}

func operatorRaftConfig(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	rc, err := client.RaftGetConfiguration(queryOpts)
	if err != nil {
		return err
	}

	return output(rc)
}
