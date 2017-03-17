package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
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

func newCatalogDatacentersCommand() *cobra.Command {
	c := action.CatalogDatacentersAction()

        cmd := &cobra.Command{
                Use:   "datacenters",
                Short: "Get all the datacenters known by the Consul server",
                Long:  "Get all the datacenters known by the Consul server",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

var deregisterLongHelp = `Deregister a service, node or check from the catalog
                
  If only --node is provided, the node and all associated services and checks are
deleted.

  If --check-id is provided, only that check is removed.

  If --service-id is provided, only that service is removed.
`

func newCatalogDeregisterCommand() *cobra.Command {
	c := action.CatalogDeregisterAction()

        cmd := &cobra.Command{
                Use:   "deregister",
                Short: "Deregisters a node, service or check",
                Long:  deregisterLongHelp,
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCatalogNodeCommand() *cobra.Command {
	c := action.CatalogNodeAction()

        cmd := &cobra.Command{
                Use:   "node",
                Short: "Get the services provided by a node",
                Long:  "Get the services provided by a node",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCatalogNodesCommand() *cobra.Command {
	c := action.CatalogNodesAction()

        cmd := &cobra.Command{
                Use:   "nodes",
                Short: "Get all the nodes registered with a given DC",
                Long:  "Get all the nodes registered with a given DC",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCatalogRegisterCommand() *cobra.Command {
	c := action.CatalogRegisterAction()

        cmd := &cobra.Command{
                Use:   "register-service <name>",
                Short: "Register external services",
                Long:  "Register external services",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCatalogServiceCommand() *cobra.Command {
	c := action.CatalogServiceAction()

        cmd := &cobra.Command{
                Use:   "service",
                Short: "Get the services provided by a service",
                Long:  "Get the services provided by a service",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCatalogServicesCommand() *cobra.Command {
	c := action.CatalogServicesAction()

        cmd := &cobra.Command{
                Use:   "services",
                Short: "Get all the services registered with a given DC",
                Long:  "Get all the services registered with a given DC",
                RunE: func(cmd *cobra.Command, args []string) error {
                       return c.Run(args)
                },
        }

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

