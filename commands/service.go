package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Consul /agent/service endpoint interface",
		Long:  "Consul /agent/service endpoint interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, []string{})
			return nil
		},
	}

	cmd.AddCommand(newServiceDeregisterCommand())
	cmd.AddCommand(newServiceMaintenanceCommand())
	cmd.AddCommand(newServiceRegisterCommand())

	return cmd
}

func newServiceDeregisterCommand() *cobra.Command {
	s := action.ServiceDeregisterAction()

        cmd := &cobra.Command{
                Use:   "deregister <serviceId>",
                Short: "Remove a service from the agent",
                Long:  "Remove a service from the agent",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

        return cmd
}

func newServiceMaintenanceCommand() *cobra.Command {
	s := action.ServiceMaintenanceAction()

        cmd := &cobra.Command{
                Use:   "maintenance",
                Short: "Manage maintenance mode of a service",
                Long:  "Manage maintenance mode of a service",
                RunE:  func (cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

        return cmd
}

var srLongHelp = `Register a new local service

  If --id is not specified, the serviceName is used. There cannot
be duplicate service IDs per agent however.

  If --address is not specified, the IP address of the local agent
is used.

  Checks are defined with the --check flag. The flags specified after
the --check flag are specific to that check. To define multiple checks,
multiple --check flags can be used:

  --check --http=http://localhost:8500/v1/agent/self --interval 30m \
  --check --http=http://localhost:8500/v1/status/leader --interval 5m --notes "Leader check"

The above example defines two checks. The first queries the /v1/agent/self endpoint
every 30 minutes. The second queries /v1/status/leader every five minutes.
`

func newServiceRegisterCommand() *cobra.Command {
	s := action.ServiceRegisterAction()

        cmd := &cobra.Command{
                Use:   "register <serviceName>",
                Short: "Register a new local service",
                Long:  srLongHelp,
                RunE:  func (cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
        }

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

        return cmd
}
