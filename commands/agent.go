package commands

import (
	"github.com/spf13/cobra"
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
