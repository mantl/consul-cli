package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Force Leave functions

func newAgentForceLeaveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force-leave <node name>",
		Short: "Force the removal of a node",
		Long:  "Force the removal of a node",
		RunE:  agentForceLeave,
	}

	return cmd
}

func agentForceLeave(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Name not provided")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.ForceLeave(args[0])
}
