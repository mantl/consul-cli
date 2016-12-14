package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Datacenters functions

func newCatalogDatacentersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Get all the datacenters known by the Consul server",
		Long:  "Get all the datacenters known by the Consul server",
		RunE:  catalogDatacenters,
	}

	addTemplateOption(cmd)

	return cmd
}

func catalogDatacenters(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	config, err := client.Datacenters()
	if err != nil {
		return err
	}

	return output(config)
}

// Deregister functions

var deregisterLongHelp = `Deregister a service, node or check from the catalog

  If only --node is provided, the node and all associated services and checks are
deleted.

  If --check-id is provided, only that check is removed.

  If --service-id is provided, only that service is removed.
`

func newCatalogDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister",
		Short: "Deregisters a node, service or check",
		Long:  deregisterLongHelp,
		RunE:  catalogDeregister,
	}

	addDatacenterOption(cmd)

	cmd.Flags().String("node", "", "Consul node name. Required")
	cmd.Flags().String("service-id", "", "Service ID to deregister")
	cmd.Flags().String("check-id", "", "Check ID to deregister")

	return cmd
}

func catalogDeregister(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	node := viper.GetString("node")
	if node == "" {
		return fmt.Errorf("Node name is required for catalog deregistration")
	}

	writeOpts := writeOptions()
	_, err = client.Deregister(&consulapi.CatalogDeregistration{
		Node:       node,
		Datacenter: viper.GetString("datacenter"),
		ServiceID:  viper.GetString("string-id"),
		CheckID:    viper.GetString("check-id"),
	}, writeOpts)

	return err
}

// Node functions

func newCatalogNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Get the services provided by a node",
		Long:  "Get the services provided by a node",
		RunE:  catalogNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogNode(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Node(args[0], queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}

// Nodes functions

func newCatalogNodesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "Get all the nodes registered with a given DC",
		Long:  "Get all the nodes registered with a given DC",
		RunE:  catalogNodes,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogNodes(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}

// Register functions

func newCatalogRegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-service <name>",
		Short: "Register external services",
		Long:  "Register external services",
		RunE:  catalogRegister,
	}

	cmd.Flags().String("node", "", "Service node")

	addDatacenterOption(cmd)
	addAgentServiceOptions(cmd)

	return cmd
}

func catalogRegister(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	if len(args) != 1 {
		return fmt.Errorf("A single service name must be specified")
	}
	service := args[0]

	client, err := newCatalog()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	_, err = client.Register(&consulapi.CatalogRegistration{
		Node:       viper.GetString("node"),
		Address:    viper.GetString("address"),
		Datacenter: viper.GetString("datacenter"),
		Service: &consulapi.AgentService{
			ID:                viper.GetString("id"),
			Service:           service,
			Tags:              getStringSlice(cmd, "tag"),
			Port:              viper.GetInt("port"),
			EnableTagOverride: viper.GetBool("override-tag"),
		},
	},
		writeOpts)

	return err
}

// Service functions

func newCatalogServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Get the services provided by a service",
		Long:  "Get the services provided by a service",
		RunE:  catalogService,
	}

	cmd.Flags().String("tag", "", "Service tag to filter on")

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogService(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Service(args[0], viper.GetString("tag"), queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}

// Services functions

func newCatalogServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Get all the services registered with a given DC",
		Long:  "Get all the services registered with a given DC",
		RunE:  catalogServices,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func catalogServices(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCatalog()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	config, _, err := client.Services(queryOpts)
	if err != nil {
		return err
	}

	return output(config)
}
