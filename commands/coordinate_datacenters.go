package commands

import (
	"github.com/spf13/cobra"
)

func (c *Coordinate) AddDatacentersSub(cmd *cobra.Command) {
	datacentersCmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Queries for WAN coordinates of Consul servers",
		Long:  "Queries for WAN coordinates of Consul servers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Datacenters(args)
		},
	}

	c.AddTemplateOption(datacentersCmd)

	cmd.AddCommand(datacentersCmd)
}

func (c *Coordinate) Datacenters(args []string) error {
	client, err := c.Coordinate()
	if err != nil {
		return err
	}

	data, err := client.Datacenters()
	if err != nil {
		return err
	}

	return c.Output(data)
}
