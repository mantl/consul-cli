package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Service functions

func newCatalogServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Get the services provided by a service",
		Long:  "Get the services provided by a service",
		RunE:  catalogService,
	}

	cmd.Flags().String("tag", "", "Service tag to filter on")

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogService(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Service(args[0], viper.GetString("tag"), queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}
