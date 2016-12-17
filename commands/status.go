package commands

import (
	"github.com/spf13/cobra"
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
