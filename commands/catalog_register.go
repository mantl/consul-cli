package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
