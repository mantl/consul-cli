package commands

import (
	"github.com/spf13/cobra"
)

func newCatalogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "catalog",
		Short: "Consul /catalog endpoint interface",
		Long:  "Consul /catalog endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newCatalogDatacentersCommand())
	cmd.AddCommand(newCatalogDeregisterCommand())
	cmd.AddCommand(newCatalogNodeCommand())
	cmd.AddCommand(newCatalogNodesCommand())
	cmd.AddCommand(newCatalogRegisterCommand())
	cmd.AddCommand(newCatalogServiceCommand())
	cmd.AddCommand(newCatalogServicesCommand())

	return cmd
}
