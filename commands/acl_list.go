package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// List functions

func newAclListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all active ACL tokens",
		Long:  "List all active ACL tokens",
		RunE:  aclList,
	}

	addTemplateOption(cmd)

	return cmd
}

func aclList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	acls, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return output(acls)
}
