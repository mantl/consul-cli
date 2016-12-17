package commands

import (
	"github.com/spf13/cobra"
)

// raft command

func newOperatorRaftCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raft",
		Short: "Consul /operator/raft endpoint interface",
		Long:  "Consul /operator/raft endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorRaftConfigCommand())
	cmd.AddCommand(newOperatorRaftDeleteCommand())

	return cmd
}
