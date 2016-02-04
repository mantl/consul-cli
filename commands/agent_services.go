package commands

import (
	"github.com/spf13/cobra"
)

func (a *Agent) AddServicesSub(c *cobra.Command) {
	servicesCmd := &cobra.Command{
		Use: "services",
		Short: "Get the services the agent is managing",
		Long: "Get the services the agent is managing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Services(args)
		},
	}

	oldServicesCmd := &cobra.Command{
		Use: "agent-services",
		Short: "Get the services the agent is managing",
		Long: "Get the services the agent is managing",
		Deprecated: "Use agent services",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Services(args)
		},
	}

	a.AddTemplateOption(servicesCmd)
	c.AddCommand(servicesCmd)

	a.AddCommand(oldServicesCmd)
}

func (a *Agent) Services(args []string) error {
	consul, err := a.Client()
	if err != nil {
		return err
	}

	client := consul.Agent()
	config, err := client.Services()
	if err != nil {
		return err
	}

	return a.Output(config)
}
