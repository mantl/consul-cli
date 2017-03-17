package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
)

func NewConsulCliCommand(name, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "consul-cli",
		Short:         "Command line interface for Consul HTTP API",
		Long:          "Command line interface for Consul HTTP API",
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, []string{})
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&cmd.SilenceUsage, "quiet", "q", true, "Don't show usage on error")
	cmd.PersistentFlags().AddGoFlagSet(action.GlobalCommandFlags())

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s %s\n", name, version)
			return nil
		},
	}
	cmd.AddCommand(versionCmd)

	cmd.AddCommand(newAclCommand())
	cmd.AddCommand(newAgentCommand())
	cmd.AddCommand(newCatalogCommand())
	cmd.AddCommand(newCheckCommand())
	cmd.AddCommand(newCoordinateCommand())
	cmd.AddCommand(newEventCommand())
	cmd.AddCommand(newHealthCommand())
	cmd.AddCommand(newKvCommand())
	cmd.AddCommand(newOperatorCommand())
	cmd.AddCommand(newServiceCommand())
	cmd.AddCommand(newSessionCommand())
	cmd.AddCommand(newStatusCommand())

	return cmd
}
