package commands

import (
	"github.com/spf13/cobra"
)

func (c *Catalog) AddServicesSub(cmd *cobra.Command) {
	servicesCmd := &cobra.Command{
		Use:   "services",
		Short: "Get all the services registered with a given DC",
		Long:  "Get all the services registered with a given DC",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Services(args)
		},
	}

	oldServicesCmd := &cobra.Command{
		Use:        "catalog-services",
		Short:      "Get all the services registered with a given DC",
		Long:       "Get all the services registered with a given DC",
		Deprecated: "Use catalog services",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Services(args)
		},
	}

	c.AddDatacenterOption(servicesCmd)
	c.AddDatacenterOption(oldServicesCmd)

	c.AddTemplateOption(servicesCmd)
	cmd.AddCommand(servicesCmd)

	c.AddCommand(oldServicesCmd)
}

func (c *Catalog) Services(args []string) error {
	client, err := c.Catalog()
	if err != nil {
		return err
	}

	queryOpts := c.QueryOptions()
	config, _, err := client.Services(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
