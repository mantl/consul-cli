package commands

import (
	"github.com/spf13/cobra"
)

func newCoordinateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coordinate",
		Short: "Consul /coordinate endpoint interface",
		Long:  "Consul /coordinate endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newCoordDatacentersCommand())
	cmd.AddCommand(newCoordNodesCommand())

	return cmd
}

