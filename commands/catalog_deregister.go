package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
