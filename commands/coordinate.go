package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Datacenters functions

func newCoordDatacentersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Queries for WAN coordinates of Consul servers",
		Long:  "Queries for WAN coordinates of Consul servers",
		RunE:  coordDatacenters,
	}

	addTemplateOption(cmd)

	return cmd
}
func coordDatacenters(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCoordinate()
	if err != nil {
		return err
	}

	data, err := client.Datacenters()
	if err != nil {
		return err
	}

	return output(data)
}

// Nodes functions

func newCoordNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Queries for LAN coordinates of Consul servers",
		Long:  "Queries for LAN coordinates of Consul servers",
		RunE:  coordNodes,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func coordNodes(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCoordinate()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	data, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return output(data)
}
