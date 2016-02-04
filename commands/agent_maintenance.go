package commands

import (
	"github.com/spf13/cobra"
)

type AgentMaintenanceOptions struct {
	enabled bool
	reason string
}

func (a *Agent) AddMaintenanceSub(c *cobra.Command) {
	amo := &AgentMaintenanceOptions{}

	maintenanceCmd := &cobra.Command{
		Use: "maintenance",
		Short: "Manage node maintenance mode",
		Long: "Manage node maintenance mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Maintenance(args, amo)
		},
	}

	oldMaintenanceCmd := &cobra.Command{
		Use: "agent-maintenance",
		Short: "Manage node maintenance mode",
		Long: "Manage node maintenance mode",
		Deprecated: "Use agent maintenance",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Maintenance(args, amo)
		},
	}

	maintenanceCmd.Flags().BoolVar(&amo.enabled, "enabled", true, "Boolean value for maintenance mode")
	maintenanceCmd.Flags().StringVar(&amo.reason, "reason", "", "Reason for entering maintenance mode")
	oldMaintenanceCmd.Flags().BoolVar(&amo.enabled, "enabled", true, "Boolean value for maintenance mode")
	oldMaintenanceCmd.Flags().StringVar(&amo.reason, "reason", "", "Reason for entering maintenance mode")

	a.AddTemplateOption(maintenanceCmd)
	c.AddCommand(maintenanceCmd)

	a.AddCommand(oldMaintenanceCmd)
}

func (a *Agent) Maintenance(args []string, amo *AgentMaintenanceOptions) error {
	consul, err := a.Client()
	if err != nil {	
		return err
	}

	client := consul.Agent()

	if amo.enabled {
		return client.EnableNodeMaintenance(amo.reason)
	} else {
		return client.DisableNodeMaintenance()
	}
}
