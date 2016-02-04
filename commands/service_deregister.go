package commands

import (
	"github.com/spf13/cobra"
)

func (s *Service) AddDeregisterSub(cmd *cobra.Command) {
	deregisterCmd := &cobra.Command{
		Use: "deregister <serviceId>",
		Short: "Remove a service from the agent",
		Long: "Remove a service from the agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Deregister(args)
		},
	}

	oldDeregisterCmd := &cobra.Command{
		Use: "service-deregister <serviceId>",
		Short: "Remove a service from the agent",
		Long: "Remove a service from the agent",
		Deprecated: "Use service deregister",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Deregister(args)
		},
	}

	cmd.AddCommand(deregisterCmd)

	s.AddCommand(oldDeregisterCmd)
}

func (s *Service) Deregister(args []string) error {
	if err := s.CheckIdArg(args); err != nil {
		return err
	}
	serviceId := args[0]

	consul, err := s.Client()
	if err != nil {	
		return err
	}

	client := consul.Agent()
	err = client.ServiceDeregister(serviceId)
	if err != nil {
		return err
	}

	return nil
}

