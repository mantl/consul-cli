package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Create functions

func newAclCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [<id>]",
		Short: "Create an ACL. Requires a management token.",
		Long:  "Create an ACL. Requires a management token.",
		RunE:  aclCreate,
	}

	cmd.Flags().Bool("management", false, "Create a management token")
	cmd.Flags().String("name", "", "Name of the ACL")
	cmd.Flags().StringSlice("rule", nil, "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	addRawOption(cmd)

	return cmd
}

func aclCreate(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = viper.GetString("name")

	switch {
	case len(args) == 1:
		entry.ID = args[0]
	case len(args) > 1:
		return fmt.Errorf("Only one ACL identified can be specified")
	}

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
	id, _, err := client.Create(entry, writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}
