package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Consul /agent endpoint interface",
		Long:  "Consul /agent endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newAgentChecksCommand())
	cmd.AddCommand(newAgentForceLeaveCommand())
	cmd.AddCommand(newAgentJoinCommand())
	cmd.AddCommand(newAgentLeaveCommand())
	cmd.AddCommand(newAgentMaintenanceCommand())
	cmd.AddCommand(newAgentMembersCommand())
	cmd.AddCommand(newAgentMonitorCommand())
	cmd.AddCommand(newAgentReloadCommand())
	cmd.AddCommand(newAgentSelfCommand())
	cmd.AddCommand(newAgentServicesCommand())

	return cmd
}

func newAgentChecksCommand() *cobra.Command {
	ag := action.AgentChecksAction()

	cmd := &cobra.Command{
                Use:   "checks",
                Short: "Get the checks the agent is managing",
                Long:  "Get the checks the agent is managing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentForceLeaveCommand() *cobra.Command {
	ag := action.AgentForceLeaveAction()

	cmd := &cobra.Command{
		                Use:   "force-leave <node name>",
                Short: "Force the removal of a node",
                Long:  "Force the removal of a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentJoinCommand() *cobra.Command {
	ag := action.AgentJoinAction()

	cmd := &cobra.Command{
                Use:   "join",
                Short: "Trigger the local agent to join a node",
                Long:  "Trigger the local agent to join a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentLeaveCommand() *cobra.Command {
	ag := action.AgentLeaveAction()

	cmd := &cobra.Command{
                Use:   "leave",
                Short: "Cause the agent to gracefully shutdown and leave the cluster",
                Long:  "Cause the agent to gracefully shutdown and leave the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentMaintenanceCommand() *cobra.Command {
	ag := action.AgentMaintenanceAction()

	cmd := &cobra.Command{
                Use:   "maintenance",
                Short: "Manage node maintenance mode",
                Long:  "Manage node maintenance mode",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentMembersCommand() *cobra.Command {
	ag := action.AgentMembersAction()

	cmd := &cobra.Command{
                Use:   "members",
                Short: "Get the members as seen by the serf agent",
                Long:  "Get the members as seen by the serf agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentMonitorCommand() *cobra.Command {
	ag := action.AgentMonitorAction()

	cmd := &cobra.Command{
                Use:   "monitor",
                Short: "Streams logs from the agent",
                Long:  "Streams logs from the agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentReloadCommand() *cobra.Command {
	ag := action.AgentReloadAction()

	cmd := &cobra.Command{
                Use:   "reload",
                Short: "Tell the Consul agent to reload its configuration",
                Long:  "Tell the Consul agent to reload its configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentSelfCommand() *cobra.Command {
	ag := action.AgentSelfAction()

	cmd := &cobra.Command{
               Use:   "self",
                Short: "Get agent configuration",
                Long:  "Get agent configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}

func newAgentServicesCommand() *cobra.Command {
	ag := action.AgentServicesAction()

	cmd := &cobra.Command{
                Use:   "services",
                Short: "Get the services the agent is managing",
                Long:  "Get the services the agent is managing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ag.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ag.CommandFlags())

	return cmd
}
