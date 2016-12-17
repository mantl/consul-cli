package commands

import (
	"github.com/spf13/cobra"
)

// keyring command

func newOperatorKeyringCommand() *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true, // Hide subcommand Consul official release
		Use:    "keyring",
		Short:  "Consul /operator/keyring interface",
		Long:   "Consul /operator/keyring interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorKeyringInstallCommand())
	cmd.AddCommand(newOperatorKeyringListCommand())
	cmd.AddCommand(newOperatorKeyringRemoveCommand())
	cmd.AddCommand(newOperatorKeyringUseCommand())

	return cmd
}
