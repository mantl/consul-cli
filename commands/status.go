package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Consul /status endpoint interface",
		Long:  "Consul /status endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newStatusLeaderCommand())
	cmd.AddCommand(newStatusPeersCommand())

	return cmd
}

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
