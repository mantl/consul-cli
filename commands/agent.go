package commands

import (
	"github.com/spf13/cobra"
)

type Agent struct {
	*Cmd
}

func (root *Cmd) initAgent() {
	a := Agent{Cmd: root}

	agentCmd := &cobra.Command{
		Use:   "agent",
		Short: "Consul /agent endpoint interface",
		Long:  "Consul /agent endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			root.Help()
		},
	}

	a.AddChecksSub(agentCmd)
	a.AddForceLeaveSub(agentCmd)
	a.AddJoinSub(agentCmd)
	a.AddMaintenanceSub(agentCmd)
	a.AddMembersSub(agentCmd)
	a.AddSelfSub(agentCmd)
	a.AddServicesSub(agentCmd)
	a.AddCommand(agentCmd)
}
