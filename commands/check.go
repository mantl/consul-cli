package commands

import (
	"github.com/spf13/cobra"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Consul /agent/check interface",
		Long:  "Consul /agent/check interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newCheckDeregisterCommand())
	cmd.AddCommand(newCheckFailCommand())
	cmd.AddCommand(newCheckPassCommand())
	cmd.AddCommand(newCheckRegisterCommand())
	cmd.AddCommand(newCheckWarnCommand())

	return cmd
}
