package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type CatalogServiceOptions struct {
	tag string
}

func (c *Catalog) AddServiceSub(cmd *cobra.Command) {
	cso := &CatalogServiceOptions{}

	serviceCmd := &cobra.Command{
		Use: "service",
		Short: "Get the services provided by a service",
		Long: "Get the services provided by a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Service(args, cso)
		},
	}

	oldServiceCmd := &cobra.Command{
		Use: "catalog-service",
		Short: "Get the services provided by a service",
		Long: "Get the services provided by a service",
		Deprecated: "Use catalog service",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Service(args, cso)
		},
	}

	serviceCmd.Flags().StringVar(&cso.tag, "tag", "", "Service tag to filter on")
	oldServiceCmd.Flags().StringVar(&cso.tag, "tag", "", "Service tag to filter on")

	c.AddDatacenterOption(serviceCmd)
	c.AddDatacenterOption(oldServiceCmd)

	c.AddTemplateOption(serviceCmd)
	cmd.AddCommand(serviceCmd)

	c.AddCommand(oldServiceCmd)
}

func (c *Catalog) Service(args []string, cso *CatalogServiceOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service allowed")
	}
	
	consul, err := c.Client()
	if err != nil {
		return err
	}

	client := consul.Catalog()
	queryOpts := c.QueryOptions()
	config, _, err := client.Service(args[0], cso.tag, queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
