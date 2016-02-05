package commands

import (
	"github.com/spf13/cobra"
)

type ServiceMaintenanceOptions struct {
	enabled bool
	reason  string
}

func (s *Service) AddMaintenanceSub(cmd *cobra.Command) {
	smo := &ServiceMaintenanceOptions{}

	maintenanceCmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage maintenance mode of a service",
		Long:  "Manage maintenance mode of a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Maintenance(args, smo)
		},
	}

	oldMaintenanceCmd := &cobra.Command{
		Use:        "service-maintenance",
		Short:      "Manage maintenance mode of a service",
		Long:       "Manage maintenance mode of a service",
		Deprecated: "Use agent maintenance",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Maintenance(args, smo)
		},
	}

	maintenanceCmd.Flags().BoolVar(&smo.enabled, "enabled", true, "Boolean value for maintenance mode")
	maintenanceCmd.Flags().StringVar(&smo.reason, "reason", "", "Reason for entering maintenance mode")
	oldMaintenanceCmd.Flags().BoolVar(&smo.enabled, "enabled", true, "Boolean value for maintenance mode")
	oldMaintenanceCmd.Flags().StringVar(&smo.reason, "reason", "", "Reason for entering maintenance mode")

	s.AddTemplateOption(maintenanceCmd)
	cmd.AddCommand(maintenanceCmd)

	s.AddCommand(oldMaintenanceCmd)
}

func (s *Service) Maintenance(args []string, smo *ServiceMaintenanceOptions) error {
	if err := s.CheckIdArg(args); err != nil {
		return err
	}
	serviceId := args[0]

	consul, err := s.Client()
	if err != nil {
		return err
	}

	client := consul.Agent()

	if smo.enabled {
		return client.EnableServiceMaintenance(serviceId, smo.reason)
	} else {
		return client.DisableServiceMaintenance(serviceId)
	}
}
