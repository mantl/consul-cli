package commands

import (
	"github.com/spf13/cobra"
)

func (c *Catalog) AddNodesSub(cmd *cobra.Command) {
	nodesCmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get all the nodes registered with a given DC",
		Long:  "Get all the nodes registered with a given DC",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Nodes(args)
		},
	}

	oldNodesCmd := &cobra.Command{
		Use:        "catalog-nodes",
		Short:      "Get all the nodes registered with a given DC",
		Long:       "Get all the nodes registered with a given DC",
		Deprecated: "Use catalog nodes",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Nodes(args)
		},
	}

	c.AddDatacenterOption(nodesCmd)
	c.AddDatacenterOption(oldNodesCmd)

	c.AddTemplateOption(nodesCmd)
	cmd.AddCommand(nodesCmd)

	c.AddCommand(oldNodesCmd)
}

func (c *Catalog) Nodes(args []string) error {
	client, err := c.Catalog()
	if err != nil {
		return err
	}

	queryOpts := c.QueryOptions()
	config, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
