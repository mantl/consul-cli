package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Deregistration functions

func newServiceDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister <serviceId>",
		Short: "Remove a service from the agent",
		Long:  "Remove a service from the agent",
		RunE:  serviceDeregister,
	}

	return cmd
}

func serviceDeregister(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	viper.BindPFlags(cmd.Flags())

	agent, err := newAgent()
	if err != nil {
		return err
	}

	var result error

	for _, id := range args {
		if err := agent.ServiceDeregister(id); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
