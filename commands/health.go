package commands

import (
	"github.com/spf13/cobra"
)

type Health struct {
	*Cmd
}

func (root *Cmd) initHealth() {
	h := Health{ Cmd: root }

	healthCmd := &cobra.Command{
		Use: "health",
		Short: "Consul Health endpoint interface",
		Long: "Consul Health endpoint interface",
		Run: func (cmd *cobra.Command, args []string) {
			root.Help()
		},
	}

	h.AddChecksSub(healthCmd)
	h.AddNodeSub(healthCmd)
	h.AddServiceSub(healthCmd)
	h.AddStateSub(healthCmd)

	h.AddCommand(healthCmd)
}

