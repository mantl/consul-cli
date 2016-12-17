package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Replication functions

func newAclReplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "replication",
		Short:  "Get the status of the ACL replication process",
		Long:   "Get the status of the ACL replication process",
		RunE:   aclReplication,
		Hidden: true,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func aclReplication(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	return fmt.Errorf("ACL replication status not available in Consul API")

	//
	//	client, err := newACL()
	//	if err != nil {
	//		return err
	//	}
}
