package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Datacenters functions

func newCoordDatacentersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datacenters",
		Short: "Queries for WAN coordinates of Consul servers",
		Long:  "Queries for WAN coordinates of Consul servers",
		RunE:  coordDatacenters,
	}

	addTemplateOption(cmd)

	return cmd
}
func coordDatacenters(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newCoordinate()
	if err != nil {
		return err
	}

	data, err := client.Datacenters()
	if err != nil {
		return err
	}

	return output(data)
}
