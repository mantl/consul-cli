package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Join functions

func newAgentJoinCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Trigger the local agent to join a node",
		Long:  "Trigger the local agent to join a node",
		RunE:  agentJoin,
	}

	cmd.Flags().Bool("wan", false, "Get list of WAN join instead of LAN")

	return cmd
}

func agentJoin(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one name allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Join(args[0], viper.GetBool("wan"))
}
