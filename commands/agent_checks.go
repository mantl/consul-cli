package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Checks functions

func newAgentChecksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks",
		Short: "Get the checks the agent is managing",
		Long:  "Get the checks the agent is managing",
		RunE:  agentChecks,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentChecks(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Checks()
	if err != nil {
		return err
	}

	return output(config)
}
