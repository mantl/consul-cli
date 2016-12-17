package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Checks functions

func newHealthChecksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks <serviceName>",
		Short: "Get the health checks for a service",
		Long:  "Get the health checks for a service",
		RunE:  healthChecks,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthChecks(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	checks, _, err := client.Checks(service, queryOpts)
	if err != nil {
		return err
	}

	return output(checks)
}
