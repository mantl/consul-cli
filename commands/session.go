package commands

import (
	"github.com/spf13/cobra"
)

func newSessionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Consul /session endpoint interface",
		Long:  "Consul /session endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newSessionCreateCommand())
	cmd.AddCommand(newSessionDestroyCommand())
	cmd.AddCommand(newSessionInfoCommand())
	cmd.AddCommand(newSessionListCommand())
	cmd.AddCommand(newSessionNodeCommand())
	cmd.AddCommand(newSessionRenewCommand())

	return cmd
}
