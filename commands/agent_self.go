package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Self functions

func newAgentSelfCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self",
		Short: "Get agent configuration",
		Long:  "Get agent configuration",
		RunE:  agentSelf,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentSelf(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Self()
	if err != nil {
		return err
	}

	return output(config)
}
