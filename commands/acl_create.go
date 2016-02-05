package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type AclCreateOptions struct {
	IsManagement bool
	Name         string
	ConfigRules  []*ConfigRule
}

func (a *Acl) AddCreateSub(c *cobra.Command) {
	aco := &AclCreateOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create an ACL. Requires a management token.",
		Long:  "Create an ACL. Requires a management token.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Create(args, aco)
		},
	}

	oldCreateCmd := &cobra.Command{
		Use:        "acl-create",
		Short:      "Create an ACL. Requires a management token.",
		Long:       "Create an ACL. Requires a management token.",
		Deprecated: "Use acl create",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Create(args, aco)
		},
	}

	createCmd.Flags().BoolVar(&aco.IsManagement, "management", false, "Create a management token")
	createCmd.Flags().StringVar(&aco.Name, "name", "", "Name of the ACL")
	createCmd.Flags().Var((funcVar)(func(s string) error {
		t, err := a.ParseRuleConfig(s)
		if err != nil {
			return err
		}

		if aco.ConfigRules == nil {
			aco.ConfigRules = make([]*ConfigRule, 0, 1)
		}

		aco.ConfigRules = append(aco.ConfigRules, t)
		return nil
	}), "rule", "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	oldCreateCmd.Flags().BoolVar(&aco.IsManagement, "management", false, "Create a management token")
	oldCreateCmd.Flags().StringVar(&aco.Name, "name", "", "Name of the ACL")
	oldCreateCmd.Flags().Var((funcVar)(func(s string) error {
		t, err := a.ParseRuleConfig(s)
		if err != nil {
			return err
		}

		if aco.ConfigRules == nil {
			aco.ConfigRules = make([]*ConfigRule, 0, 1)
		}

		aco.ConfigRules = append(aco.ConfigRules, t)
		return nil
	}), "rule", "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	c.AddCommand(createCmd)

	a.AddCommand(oldCreateCmd)
}

func (a *Acl) Create(args []string, aco *AclCreateOptions) error {
	client, err := a.ACL()
	if err != nil {
		return err
	}

	var entry *consulapi.ACLEntry

	if aco.IsManagement {
		entry = &consulapi.ACLEntry{
			Name: aco.Name,
			Type: consulapi.ACLManagementType,
		}
	} else {
		rules, err := a.GetRulesString(aco.ConfigRules)
		if err != nil {
			return err
		}

		entry = &consulapi.ACLEntry{
			Name:  aco.Name,
			Type:  consulapi.ACLClientType,
			Rules: rules,
		}

	}

	writeOpts := a.WriteOptions()
	id, _, err := client.Create(entry, writeOpts)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, id)

	return nil
}
