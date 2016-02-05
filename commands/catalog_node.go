package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *Catalog) AddNodeSub(cmd *cobra.Command) {
	nodeCmd := &cobra.Command{
		Use:   "node",
		Short: "Get the services provided by a node",
		Long:  "Get the services provided by a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Node(args)
		},
	}

	oldNodeCmd := &cobra.Command{
		Use:        "catalog-node",
		Short:      "Get the services provided by a node",
		Long:       "Get the services provided by a node",
		Deprecated: "Use catalog node",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Node(args)
		},
	}

	c.AddDatacenterOption(nodeCmd)
	c.AddDatacenterOption(oldNodeCmd)

	c.AddTemplateOption(nodeCmd)
	cmd.AddCommand(nodeCmd)

	c.AddCommand(oldNodeCmd)
}

func (c *Catalog) Node(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	client, err := c.Catalog()
	if err != nil {
		return err
	}

	queryOpts := c.QueryOptions()
	config, _, err := client.Node(args[0], queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
