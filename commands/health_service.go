package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Service functions

func newHealthServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service <serviceName>",
		Short: "Get the nodes and health info for a service",
		Long:  "Get the nodes and health info for a service",
		RunE:  healthService,
	}

	cmd.Flags().String("tag", "", "Service tag to filter on")
	cmd.Flags().Bool("passing", false, "Only return passing checks")
	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthService(cmd *cobra.Command, args []string) error {
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

	tag := viper.GetString("tag")
	po := viper.GetBool("passing")

	s, _, err := client.Service(service, tag, po, queryOpts)
	if err != nil {
		return err
	}

	return output(s)
}

