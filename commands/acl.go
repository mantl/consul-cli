package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newAclCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Consul /acl endpoint interface",
		Long:  "Consul /acl endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newAclCloneCommand())
	cmd.AddCommand(newAclCreateCommand())
	cmd.AddCommand(newAclDestroyCommand())
	cmd.AddCommand(newAclInfoCommand())
	cmd.AddCommand(newAclListCommand())
	cmd.AddCommand(newAclReplicationCommand())
	cmd.AddCommand(newAclUpdateCommand())

	return cmd
}

func newAclCloneCommand() *cobra.Command {
	ac := action.AclCloneAction()

	cmd := &cobra.Command{
		Use:   "clone <id>",
		Short: "Create a new token from an existing one",
		Long:  "Create a new token from an existing one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclCreateCommand() *cobra.Command {
	ac := action.AclCreateAction()

	cmd := &cobra.Command{
		Use:   "create [<id>]",
		Short: "Create an ACL. Requires a management token.",
		Long:  "Create an ACL. Requires a management token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclDestroyCommand() *cobra.Command {
	ac := action.AclDestroyAction()

	cmd := &cobra.Command{
                Use:   "destroy <token>",
                Short: "Destroy an ACL",
                Long:  "Destroy an ACL",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclInfoCommand() *cobra.Command {
	ac := action.AclInfoAction()

	cmd := &cobra.Command{
                Use:   "info <token>",
                Short: "Query information about an ACL token",
                Long:  "Query information about an ACL token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclListCommand() *cobra.Command {
	ac := action.AclListAction()

	cmd := &cobra.Command{
                Use:   "list",
                Short: "List all active ACL tokens",
                Long:  "List all active ACL tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclReplicationCommand() *cobra.Command {
	ac := action.AclReplicationAction()

	cmd := &cobra.Command{
                Use:    "replication",
                Short:  "Get the status of the ACL replication process",
                Long:   "Get the status of the ACL replication process",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

func newAclUpdateCommand() *cobra.Command {
	ac := action.AclUpdateAction()

	cmd := &cobra.Command{
                Use:   "update <token>",
                Short: "Update an ACL. Will be created if it doesn't exist",
                Long:  "Update an ACL. Will be created if it doesn't exist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ac.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(ac.CommandFlags())

	return cmd
}

