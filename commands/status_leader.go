package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Leader functions

func newStatusLeaderCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leader",
		Short: "Get the current Raft leader",
		Long:  "Get the current Raft leader",
		RunE:  statusLeader,
	}

	addTemplateOption(cmd)

	return cmd
}

func statusLeader(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newStatus()
	if err != nil {
		return err
	}

	l, err := client.Leader()
	if err != nil {
		return err
	}

	return output(l)
}
