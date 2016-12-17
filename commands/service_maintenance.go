package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Maintenance functions

func newServiceMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage maintenance mode of a service",
		Long:  "Manage maintenance mode of a service",
		RunE:  serviceMaintenance,
	}

	cmd.Flags().Bool("enabled", true, "Boolean value for maintenance mode")
	cmd.Flags().String("reason", "", "Reason for entering maintenance mode")

	return cmd
}

func serviceMaintenance(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	viper.BindPFlags(cmd.Flags())

	agent, err := newAgent()
	if err != nil {
		return err
	}

	var result error

	enabled := viper.GetBool("enabled")
	reason := viper.GetString("reason")

	for _, id := range args {
		if enabled {
			if err := agent.EnableServiceMaintenance(id, reason); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			if err := agent.DisableServiceMaintenance(id); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return result
}
