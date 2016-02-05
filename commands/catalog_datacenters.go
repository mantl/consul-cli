package commands

import (
	"github.com/spf13/cobra"
)

func (c *Catalog) AddDatacentersSub(cmd *cobra.Command) {
	datacentersCmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Get all the datacenters known by the Consul server",
		Long:  "Get all the datacenters known by the Consul server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Datacenters(args)
		},
	}

	oldDatacentersCmd := &cobra.Command{
		Use:        "catalog-datacenters",
		Short:      "Get all the datacenters known by the Consul server",
		Long:       "Get all the datacenters known by the Consul server",
		Deprecated: "Use catalog datacenters",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Datacenters(args)
		},
	}

	c.AddTemplateOption(datacentersCmd)
	cmd.AddCommand(datacentersCmd)

	c.AddCommand(oldDatacentersCmd)
}

func (c *Catalog) Datacenters(args []string) error {
	client, err := c.Catalog()
	if err != nil {
		return err
	}

	config, err := client.Datacenters()
	if err != nil {
		return err
	}

	return c.Output(config)
}
