package commands

import (
	"github.com/spf13/cobra"
)

func newEventCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Consul /event endpoint interface",
		Long:  "Consul /event endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newEventFireCommand())
	cmd.AddCommand(newEventListCommand())

	return cmd
}

