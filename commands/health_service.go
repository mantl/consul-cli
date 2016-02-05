package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type HealthServiceOptions struct {
	Tag         string
	PassingOnly bool
}

func (h *Health) AddServiceSub(cmd *cobra.Command) {
	hso := &HealthServiceOptions{}

	serviceCmd := &cobra.Command{
		Use:   "service <serviceName>",
		Short: "Get the nodes and health info for a service",
		Long:  "Get the nodes and health info for a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Service(args, hso)
		},
	}

	oldServiceCmd := &cobra.Command{
		Use:        "health-service <serviceName>",
		Short:      "Get the nodes and health info for a service",
		Long:       "Get the nodes and health info for a service",
		Deprecated: "Use health service",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Service(args, hso)
		},
	}

	serviceCmd.Flags().StringVar(&hso.Tag, "tag", "", "Service tag to filter on")
	serviceCmd.Flags().BoolVar(&hso.PassingOnly, "passing", false, "Only return passing checks")
	h.AddDatacenterOption(serviceCmd)
	h.AddTemplateOption(serviceCmd)

	oldServiceCmd.Flags().StringVar(&hso.Tag, "tag", "", "Service tag to filter on")
	oldServiceCmd.Flags().BoolVar(&hso.PassingOnly, "passing", false, "Only return passing checks")
	h.AddDatacenterOption(oldServiceCmd)

	cmd.AddCommand(serviceCmd)

	h.AddCommand(oldServiceCmd)
}

func (h *Health) Service(args []string, hso *HealthServiceOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	client, err := h.Health()
	if err != nil {
		return err
	}

	queryOpts := h.QueryOptions()

	s, _, err := client.Service(service, hso.Tag, hso.PassingOnly, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(s)
}
