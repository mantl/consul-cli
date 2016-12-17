package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Update functions

func newAclUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <token>",
		Short: "Update an ACL. Will be created if it doesn't exist",
		Long:  "Update an ACL. Will be created if it doesn't exist",
		RunE:  aclUpdate,
	}

	cmd.Flags().Bool("management", false, "Create a management token")
	cmd.Flags().String("name", "", "Name of the ACL")
	cmd.Flags().StringSlice("rule", nil, "Rule to update. Can be multiple rules on a command line. Format is type:path:policy")

	addRawOption(cmd)

	return cmd
}

func aclUpdate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = viper.GetString("name")
	entry.ID = id

	if viper.GetBool("management") {
		entry.Type = consulapi.ACLManagementType
	} else {
		entry.Type = consulapi.ACLClientType
	}

	if raw := viper.GetString("raw"); raw != "" {
		rules, err := readRawString(raw)
		if err != nil {
			return err
		}
		entry.Rules = rules
	} else {
		rules, err := getRulesString(getStringSlice(cmd, "rule"))
		if err != nil {
			return err
		}
		entry.Rules = rules
	}

	writeOpts := writeOptions()
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
