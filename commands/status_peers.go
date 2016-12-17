package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Peers functions

func newStatusPeersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "peers",
		Short: "Get the current Raft peers",
		Long:  "Get the current Raft peers",
		RunE:  statusPeers,
	}

	addTemplateOption(cmd)

	return cmd
}

func statusPeers(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newStatus()
	if err != nil {
		return err
	}

	l, err := client.Peers()
	if err != nil {
		return err
	}

	return output(l)
}
