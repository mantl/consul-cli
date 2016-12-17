package commands

import (
	"github.com/spf13/cobra"
)

func newOperatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator",
		Short: "Consul /operator endpoint interface",
		Long:  "Consul /operator endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorKeyringCommand())
	cmd.AddCommand(newOperatorRaftCommand())

	return cmd
}
