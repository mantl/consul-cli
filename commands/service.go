package commands

import (
	"github.com/spf13/cobra"
)

func newServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Consul /agent/service endpoint interface",
		Long:  "Consul /agent/service endpoint interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, []string{})
			return nil
		},
	}

	cmd.AddCommand(newServiceDeregisterCommand())
	cmd.AddCommand(newServiceMaintenanceCommand())
	cmd.AddCommand(newServiceRegisterCommand())

	return cmd
}
