package commands

import (
	"github.com/spf13/cobra"
	consulapi "github.com/hashicorp/consul/api"
)

type AclUpdateOptions struct {
	IsManagement bool
	Name string
	ConfigRules []*ConfigRule
}

func (a *Acl) AddUpdateSub(c *cobra.Command) {
	auo := &AclUpdateOptions{}

	updateCmd := &cobra.Command{
		Use: "update <token>",
		Short: "Update an ACL. Will be created if it doesn't exist",
		Long: "Update an ACL. Will be created if it doesn't exist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Update(args, auo)
		},
	}

	oldUpdateCmd := &cobra.Command{
		Use: "acl-update <token>",
		Short: "Update an ACL. Will be created if it doesn't exist",
		Long: "Update an ACL. Will be created if it doesn't exist",
                Deprecated: "Use acl update",
                Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Update(args, auo)
		},
	}

        updateCmd.Flags().BoolVar(&auo.IsManagement, "management", false, "Create a management token")
        updateCmd.Flags().StringVar(&auo.Name, "name", "", "Name of the ACL")
        updateCmd.Flags().Var((funcVar)(func(s string) error {
                t, err := a.ParseRuleConfig(s)
                if err != nil {
                        return err
                }

                if auo.ConfigRules == nil {
                        auo.ConfigRules = make([]*ConfigRule, 0, 1)
                }

                auo.ConfigRules = append(auo.ConfigRules, t)
                return nil
        }), "rule", "")

        oldUpdateCmd.Flags().BoolVar(&auo.IsManagement, "management", false, "Create a management token")
        oldUpdateCmd.Flags().StringVar(&auo.Name, "name", "", "Name of the ACL")
        oldUpdateCmd.Flags().Var((funcVar)(func(s string) error {
                t, err := a.ParseRuleConfig(s)
                if err != nil {
                        return err
                }

                if auo.ConfigRules == nil {
                        auo.ConfigRules = make([]*ConfigRule, 0, 1)
                }

                auo.ConfigRules = append(auo.ConfigRules, t)
                return nil
        }), "rule", "")

	c.AddCommand(updateCmd)

	a.AddCommand(oldUpdateCmd)
}

func (a *Acl) Update(args []string, auo *AclUpdateOptions) error {
	if err := a.CheckIdArg(args); err != nil {
		return err
	}
	id := args[0]

	consul, err := a.Client()
	if err != nil {
		return err
	}
	client := consul.ACL()

	var entry *consulapi.ACLEntry

	if auo.IsManagement {
		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	auo.Name,
			Type:	consulapi.ACLManagementType,
		}
	} else {
		rules, err := a.GetRulesString(auo.ConfigRules)
		if err != nil {
			return err
		}

		entry = &consulapi.ACLEntry{
			ID:	id,
			Name:	auo.Name,
			Type:	consulapi.ACLClientType,
			Rules:	rules,
		}

	}

	writeOpts := a.WriteOptions()
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
