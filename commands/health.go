package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
)

func newHealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Consul /health endpoint interface",
		Long:  "Consul /health endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newHealthChecksCommand())
	cmd.AddCommand(newHealthNodeCommand())
	cmd.AddCommand(newHealthServiceCommand())
	cmd.AddCommand(newHealthStateCommand())

	return cmd
}

func newHealthChecksCommand() *cobra.Command {
	h := action.HealthChecksAction()

	cmd := &cobra.Command{
                Use:   "checks <serviceName>",
                Short: "Get the health checks for a service",
                Long:  "Get the health checks for a service",
		RunE: func(cmd *cobra.Command, args[]string) error {
			return h.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(h.CommandFlags())

	return cmd
}

func newHealthNodeCommand() *cobra.Command {
	h := action.HealthNodeAction()

	cmd := &cobra.Command{
                Use:   "node <nodeName>",
                Short: "Get the health info for a node",
                Long:  "Get the health info for a node",
		RunE: func(cmd *cobra.Command, args[]string) error {
			return h.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(h.CommandFlags())

	return cmd
}

func newHealthServiceCommand() *cobra.Command {
	h := action.HealthServiceAction()

	cmd := &cobra.Command{
                Use:   "service <serviceName>",
                Short: "Get the nodes and health info for a service",
                Long:  "Get the nodes and health info for a service",
		RunE: func(cmd *cobra.Command, args[]string) error {
			return h.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(h.CommandFlags())

	return cmd
}

func newHealthStateCommand() *cobra.Command {
	h := action.HealthStateAction()

	cmd := &cobra.Command{
                Use:   "state",
                Short: "Get the checks in a given state",
                Long:  "Get the checks in a given state",
		RunE: func(cmd *cobra.Command, args[]string) error {
			return h.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(h.CommandFlags())

	return cmd
}

