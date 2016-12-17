package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Info functions

func newAclInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <token>",
		Short: "Query information about an ACL token",
		Long:  "Query information about an ACL token",
		RunE:  aclInfo,
	}

	addTemplateOption(cmd)

	return cmd
}

func aclInfo(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	acl, _, err := client.Info(id, queryOpts)
	if err != nil {
		return err
	}

	return output(acl)
}
