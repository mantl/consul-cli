package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Service struct {
	*Cmd
}

func (root *Cmd) initService() {
	s := Service{ Cmd: root }

	serviceCmd := &cobra.Command{
		Use: "service",
		Short: "Consul Service endpoint interface",
		Long: "Consul Service endpoint interface",
		Run: func (cmd *cobra.Command, args []string) {
			root.Help()
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
