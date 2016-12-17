package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Node functions

func newCatalogNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Get the services provided by a node",
		Long:  "Get the services provided by a node",
		RunE:  catalogNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogNode(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Node(args[0], queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}
