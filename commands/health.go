package commands

import (
	"github.com/spf13/cobra"
)

type Health struct {
	*Cmd
}

func (root *Cmd) initHealth() {
	h := Health{Cmd: root}

	healthCmd := &cobra.Command{
		Use:   "health",
		Short: "Consul /health endpoint interface",
		Long:  "Consul /health endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	h.AddChecksSub(healthCmd)
	h.AddNodeSub(healthCmd)
	h.AddServiceSub(healthCmd)
	h.AddStateSub(healthCmd)

	h.AddCommand(healthCmd)
}
