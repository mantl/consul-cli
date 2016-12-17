package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Services functions

func newCatalogServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Get all the services registered with a given DC",
		Long:  "Get all the services registered with a given DC",
		RunE:  catalogServices,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogServices(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Services(queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}
