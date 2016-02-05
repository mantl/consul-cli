package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Session struct {
	*Cmd
}

func (root *Cmd) initSession() {
	s := Session{Cmd: root}

	sessionCmd := &cobra.Command{
		Use:   "session",
		Short: "Consul /session endpoint interface",
		Long:  "Consul /session endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			root.Help()
		},
	}

	s.AddCreateSub(sessionCmd)
	s.AddDestroySub(sessionCmd)
	s.AddInfoSub(sessionCmd)
	s.AddListSub(sessionCmd)
	s.AddNodeSub(sessionCmd)
	s.AddRenewSub(sessionCmd)

	s.AddCommand(sessionCmd)
}

func (s *Session) CheckIdArg(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service id must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service id allowed")
	}

	return nil
}
