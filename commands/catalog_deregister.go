package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

func (c *Catalog) AddDeregisterSub(cmd *cobra.Command) {
	deregisterCmd := &cobra.Command{
		Use:   "deregister <serviceId> <nodeName>",
		Short: "Remove entry from the catalog",
		Long:  "Remove entry from the catalog",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Deregister(args)
		},
	}

	cmd.AddCommand(deregisterCmd)
}

func (c *Catalog) Deregister(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("No service id specified")
	case len(args) == 1:
		return fmt.Errorf("No node name specified")
	case len(args) > 2:
		return fmt.Errorf("Only pair of service id / node name is allowed")
	}

	cdo := &consulapi.CatalogDeregistration{
		ServiceID: args[0],
		Node:      args[1],
	}

	client, err := c.Catalog()
	if err != nil {
		return err
	}

	writeOpts := c.WriteOptions()

	_, err = client.Deregister(cdo, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
