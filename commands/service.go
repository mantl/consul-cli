package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Service struct {
	*Cmd
}

func (root *Cmd) initService() {
	s := Service{Cmd: root}

	serviceCmd := &cobra.Command{
		Use:   "service",
		Short: "Consul /agent/service endpoint interface",
		Long:  "Consul /agent/service endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	s.AddDeregisterSub(serviceCmd)
	s.AddMaintenanceSub(serviceCmd)
	s.AddRegisterSub(serviceCmd)

	s.AddCommand(serviceCmd)
}

func (s *Service) CheckIdArg(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("No service id specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service id allowed")
	}

	return nil
}
