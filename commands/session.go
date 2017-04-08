package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newSessionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Consul /session endpoint interface",
		Long:  "Consul /session endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newSessionCreateCommand())
	cmd.AddCommand(newSessionDestroyCommand())
	cmd.AddCommand(newSessionInfoCommand())
	cmd.AddCommand(newSessionListCommand())
	cmd.AddCommand(newSessionNodeCommand())
	cmd.AddCommand(newSessionRenewCommand())

	return cmd
}

func newSessionCreateCommand() *cobra.Command {
	s := action.SessionCreateAction()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new session",
		Long:  "Create a new session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSessionDestroyCommand() *cobra.Command {
	s := action.SessionDestroyAction()

	cmd := &cobra.Command{
		Use:   "destroy <sessionId>",
		Short: "Destroy a session",
		Long:  "Destroy a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSessionInfoCommand() *cobra.Command {
	s := action.SessionInfoAction()

	cmd := &cobra.Command{
		Use:   "info <sessionId>",
		Short: "Get information on a session",
		Long:  "Get information on a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSessionListCommand() *cobra.Command {
	s := action.SessionListAction()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active sessions for a datacenter",
		Long:  "List active sessions for a datacenter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSessionNodeCommand() *cobra.Command {
	s := action.SessionNodeAction()

	cmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get active sessions for a node",
		Long:  "Get active sessions for a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSessionRenewCommand() *cobra.Command {
	s := action.SessionRenewAction()

	cmd := &cobra.Command{
		Use:   "renew <sessionId>",
		Short: "Renew the given session",
		Long:  "Renew the given session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}
