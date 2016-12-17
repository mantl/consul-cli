package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Reload functions

func newAgentReloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Tell the Consul agent to reload its configuration",
		Long:  "Tell the Consul agent to reload its configuration",
		RunE:  agentReload,
	}

	return cmd
}

func agentReload(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Reload()
}
