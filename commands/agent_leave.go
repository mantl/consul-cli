package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Leave functions

func newAgentLeaveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave",
		Short: "Cause the agent to gracefully shutdown and leave the cluster",
		Long:  "Cause the agent to gracefully shutdown and leave the cluster",
		RunE:  agentLeave,
	}

	return cmd
}

func agentLeave(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Leave()
}
