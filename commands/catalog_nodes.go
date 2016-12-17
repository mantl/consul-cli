package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Nodes functions

func newCatalogNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get all the nodes registered with a given DC",
		Long:  "Get all the nodes registered with a given DC",
		RunE:  catalogNodes,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogNodes(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}

