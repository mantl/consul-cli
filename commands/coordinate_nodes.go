package commands

import (
	"github.com/spf13/cobra"
)

func (c *Coordinate) AddNodesSub(cmd *cobra.Command) {
	nodesCmd := &cobra.Command{
		Use:   "nodes",
		Short: "Queries for LAN coordinates of Consul servers",
		Long:  "Queries for LAN coordinates of Consul servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Nodes(args)
		},
	}

	c.AddDatacenterOption(nodesCmd)
	c.AddTemplateOption(nodesCmd)
	c.AddConsistency(nodesCmd)

	cmd.AddCommand(nodesCmd)
}

func (c *Coordinate) Nodes(args []string) error {
	client, err := c.Coordinate()
	if err != nil {
		return err
	}

	queryOpts := c.QueryOptions()
	data, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(data)
}
