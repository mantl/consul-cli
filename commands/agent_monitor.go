package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Monitor functions

func newAgentMonitorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Streams logs from the agent",
		Long:  "Streams logs from the agent",
		RunE:  agentMonitor,
	}

	cmd.Flags().String("loglevel", "", "Log level to filter on. Default is info")

	return cmd
}

func agentMonitor(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	outputChan, err := client.Monitor(
		viper.GetString("loglevel"),
		nil, // XXX - Set up interrupts and stop things gracefully
		nil, // query options. No documentation on what they are used for. Leave nil for now.
	)
	if err != nil {
		return err
	}

	for s := range outputChan {
		fmt.Printf(s)
	}

	return nil
}
