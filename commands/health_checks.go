package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)


func (h *Health) AddChecksSub(cmd *cobra.Command) {
	checksCmd := &cobra.Command{
		Use: "checks <serviceName>",
		Short: "Get the health checks for a service",
		Long: "Get the health checks for a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Checks(args)
		},
	}

	oldChecksCmd := &cobra.Command{
		Use: "health-checks <serviceName>",
		Short: "Get the health checks for a service",
		Long: "Get the health checks for a service",
		Deprecated: "Use health checks",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Checks(args)
		},
	}

	h.AddDatacenterOption(checksCmd)
	h.AddTemplateOption(checksCmd)
	h.AddDatacenterOption(oldChecksCmd)

	cmd.AddCommand(checksCmd)

	h.AddCommand(oldChecksCmd)
}

func (h *Health) Checks(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	client, err := h.Client()
	if err != nil {
		return err
	}

	queryOpts := h.QueryOptions()
	healthClient := client.Health()

	checks, _, err := healthClient.Checks(service, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(checks)
}

