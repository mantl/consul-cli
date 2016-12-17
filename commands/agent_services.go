package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Services functions

func newAgentServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Get the services the agent is managing",
		Long:  "Get the services the agent is managing",
		RunE:  agentServices,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentServices(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Services()
	if err != nil {
		return err
	}

	return output(config)
}
