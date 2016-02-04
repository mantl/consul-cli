package commands

import (
	"github.com/spf13/cobra"
)

type Status struct {
	*Cmd
}

func (root *Cmd) initStatus() {
	s := Status{ Cmd: root }

	statusCmd := &cobra.Command{
		Use: "status",
		Short: "Consul Status endpoint interface",
		Long: "Consul Status endpoint interface",
		Run: func (cmd *cobra.Command, args []string) {
			root.Help()
		},
	}

	s.AddLeaderSub(statusCmd)
	s.AddPeersSub(statusCmd)

	s.AddCommand(statusCmd)
}

