package commands

import (
	"github.com/spf13/cobra"
)

func newHealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Consul /health endpoint interface",
		Long:  "Consul /health endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newHealthChecksCommand())
	cmd.AddCommand(newHealthNodeCommand())
	cmd.AddCommand(newHealthServiceCommand())
	cmd.AddCommand(newHealthStateCommand())

	return cmd
}

