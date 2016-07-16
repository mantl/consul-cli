package commands

import (
	"github.com/spf13/cobra"
)

type Status struct {
	*Cmd
}

func (root *Cmd) initStatus() {
	s := Status{Cmd: root}

	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Consul /status endpoint interface",
		Long:  "Consul /status endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	s.AddLeaderSub(statusCmd)
	s.AddPeersSub(statusCmd)

	s.AddCommand(statusCmd)
}
