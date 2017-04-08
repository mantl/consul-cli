package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
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

func newEventFireCommand() *cobra.Command {
	e := action.EventFireAction()

	cmd := &cobra.Command{
		Use:   "fire <name>",
		Short: "Fires a new user event",
		Long:  "Fires a new user event",
		RunE: func(cmd *cobra.Command, args []string) error {
			return e.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(e.CommandFlags())

	return cmd
}

func newEventListCommand() *cobra.Command {
	e := action.EventListAction()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the most recent events the agent has seen",
		Long:  "Lists the most recent events the agent has seen",
		RunE: func(cmd *cobra.Command, args []string) error {
			return e.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(e.CommandFlags())

	return cmd
}
