package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Maintenance functions

func newAgentMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage node maintenance mode",
		Long:  "Manage node maintenance mode",
		RunE:  agentMaintenance,
	}

	cmd.Flags().Bool("enabled", true, "Boolean value for maintenance mode")
	cmd.Flags().String("reason", "", "Reason for entering maintenance mode")

	return cmd
}

func agentMaintenance(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	if viper.GetBool("enabled") {
		return client.EnableNodeMaintenance(viper.GetString("reason"))
	} 

	return client.DisableNodeMaintenance()
}
