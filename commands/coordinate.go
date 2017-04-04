package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
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

func newCoordDatacentersCommand() *cobra.Command {
	c := action.CoordDatacentersAction()

        cmd := &cobra.Command{
                Use:   "datacenters",
                Short: "Queries for WAN coordinates of Consul servers",
                Long:  "Queries for WAN coordinates of Consul servers",
                RunE: func(cmd *cobra.Command, args []string) error {
                        return c.Run(args)
                },
        }

        cmd.Flags().AddGoFlagSet(c.CommandFlags())

        return cmd
}

func newCoordNodesCommand() *cobra.Command {
	c := action.CoordNodesAction()

        cmd := &cobra.Command{
                Use:   "nodes",
                Short: "Queries for LAN coordinates of Consul servers",
                Long:  "Queries for LAN coordinates of Consul servers",
                RunE: func(cmd *cobra.Command, args []string) error {
                        return c.Run(args)
                },
        }

        cmd.Flags().AddGoFlagSet(c.CommandFlags())

        return cmd
}

