package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
