package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Datacenters functions

func newCatalogDatacentersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Get all the datacenters known by the Consul server",
		Long:  "Get all the datacenters known by the Consul server",
		RunE:  catalogDatacenters,
	}

	addTemplateOption(cmd)

	return cmd
}

func catalogDatacenters(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	config, err := client.Datacenters()
	if err != nil {
		return err
	}

	return output(config)
}
