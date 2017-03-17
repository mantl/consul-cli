package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
)

func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Consul /status endpoint interface",
		Long:  "Consul /status endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newStatusLeaderCommand())
	cmd.AddCommand(newStatusPeersCommand())

	return cmd
}

func newStatusLeaderCommand() *cobra.Command {
	s := action.StatusLeaderAction()

        cmd := &cobra.Command{
                Use:   "leader",
                Short: "Get the current Raft leader",
                Long:  "Get the current Raft leader",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newStatusPeersCommand() *cobra.Command {
	s := action.StatusPeersAction()

        cmd := &cobra.Command{
		Use:   "peers",
		Short: "Get the current Raft peers",
		Long:  "Get the current Raft peers",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}
