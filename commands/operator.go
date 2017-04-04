package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newOperatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator",
		Short: "Consul /operator endpoint interface",
		Long:  "Consul /operator endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorKeyringCommand())
	cmd.AddCommand(newOperatorRaftCommand())
//	cmd.AddCommand(newOperatorAutopilotCommand())

	return cmd
}

func newOperatorKeyringCommand() *cobra.Command {
        cmd := &cobra.Command{
                Hidden: true, // Hide subcommand Consul official release
                Use:    "keyring",
                Short:  "Consul /operator/keyring interface",
                Long:   "Consul /operator/keyring interface",
                Run: func(cmd *cobra.Command, args []string) {
                        cmd.HelpFunc()(cmd, []string{})
                },
        }

        cmd.AddCommand(newOperatorKeyringInstallCommand())
        cmd.AddCommand(newOperatorKeyringListCommand())
        cmd.AddCommand(newOperatorKeyringRemoveCommand())
        cmd.AddCommand(newOperatorKeyringUseCommand())

        return cmd
}

func newOperatorKeyringInstallCommand() *cobra.Command {
	o := action.OperatorKeyringInstallAction()

	        cmd := &cobra.Command{
                Use:   "install <key> [<key>]",
                Short: "Install a new gossip key into the cluster",
                Long:  "Install a new gossip key into the cluster",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

func newOperatorKeyringListCommand() *cobra.Command {
	o := action.OperatorKeyringListAction()

        cmd := &cobra.Command{
                Use:   "list",
                Short: "List gossip keys installed",
                Long:  "List gossip keys installed",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

func newOperatorKeyringRemoveCommand() *cobra.Command {
	o := action.OperatorKeyringRemoveAction()

        cmd := &cobra.Command{
                Use:   "remove <key> [<key>]",
		Short: "Remove gossip keys from the cluster",
		Long:  "Remove gossip keys from the cluster",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

func newOperatorKeyringUseCommand() *cobra.Command {
	o := action.OperatorKeyringUseAction()

        cmd := &cobra.Command{
                Use:   "use <key>",
                Short: "Change the primary gossip encryption key",
                Long:  "Change the primary gossip encryption key",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

func newOperatorRaftCommand() *cobra.Command {
        cmd := &cobra.Command{
                Use:   "raft",
                Short: "Consul /operator/raft endpoint interface",
                Long:  "Consul /operator/raft endpoint interface",
                Run: func(cmd *cobra.Command, args []string) {
                        cmd.HelpFunc()(cmd, []string{})
                },
        }

        cmd.AddCommand(newOperatorRaftConfigCommand())
        cmd.AddCommand(newOperatorRaftDeleteCommand())

        return cmd
}

func newOperatorRaftConfigCommand() *cobra.Command {
	o := action.OperatorRaftConfigAction()

        cmd := &cobra.Command{
                Use:   "config",
                Short: "Inspect the Raft configuration",
                Long:  "Inspect the Raft configuration",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

var raftDeleteLongHelp = `
Remove a Consul server from the Raft configuration
                
An address is required and should be set to IP:port for the server
to remove. The port number is 8300 unless configured otherwise`

func newOperatorRaftDeleteCommand() *cobra.Command {
	o := action.OperatorRaftDeleteAction()

        cmd := &cobra.Command{
                Use:   "delete <address>",
                Short: "Remove a Consul server from the Raft configuration",
                Long:  raftDeleteLongHelp,
                RunE:  func (cmd *cobra.Command, args []string) error {
			return o.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(o.CommandFlags())

        return cmd
}

//func newOperatorAutopilotCommand() *cobra.Command {
//        cmd := &cobra.Command{
//                Use:   "autopilot",
//                Short: "Consul /operator/autopilot endpoint interface",
//                Long:  "Consul /operator/autopilot endpoint interface",
//                Run: func(cmd *cobra.Command, args []string) {
//                        cmd.HelpFunc()(cmd, []string{})
//                },
//        }
//
//        cmd.AddCommand(newOperatorAutopilotGetCommand())
//        cmd.AddCommand(newOperatorAutopilotSetCommand())
//
//        return cmd
//}
//
//func newOperatorAutopilotGetCommand() *cobra.Command {
//	o := action.OperatorAutopilotGetAction()
//
//        cmd := &cobra.Command{
//                Use:   "get",
//                Short: "Retrieve the latest autopilot configuration",
//                Long: "Retrieve the latest autopilot configuration",
//                RunE:  func (cmd *cobra.Command, args []string) error {
//			return o.Run(args)
//		},
//        }
//
//	cmd.Flags().AddGoFlagSet(o.CommandFlags())
//
//        return cmd
//}
//
//func newOperatorAutopilotSetCommand() *cobra.Command {
//	o := action.OperatorAutopilotSetAction()
//
//        cmd := &cobra.Command{
//                Use:   "set",
//                Short: "Update the autopilot configuration",
//                Long: "Update the autopilot configuration",
//                RunE:  func (cmd *cobra.Command, args []string) error {
//			return o.Run(args)
//		},
//        }
//
//	cmd.Flags().AddGoFlagSet(o.CommandFlags())
//
//        return cmd
//}
