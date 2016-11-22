package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type AclUpdateOptions struct {
	IsManagement bool
	Name         string
	ConfigRules  []*ConfigRule
	Raw          string
}

func (a *Acl) AddUpdateSub(c *cobra.Command) {
	auo := &AclUpdateOptions{}

	updateCmd := &cobra.Command{
		Use:   "update <token>",
		Short: "Update an ACL. Will be created if it doesn't exist",
		Long:  "Update an ACL. Will be created if it doesn't exist",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Update(args, auo)
		},
	}

	oldUpdateCmd := &cobra.Command{
		Use:        "acl-update <token>",
		Short:      "Update an ACL. Will be created if it doesn't exist",
		Long:       "Update an ACL. Will be created if it doesn't exist",
		Deprecated: "Use acl update",
		Hidden:     true,
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
	updateCmd.Flags().StringVar(&auo.Raw, "raw", "", "Raw ACL rule definition")

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
	if !a.CheckIdArg(args) {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	client, err := a.ACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = auo.Name
	entry.ID = id

	if auo.IsManagement {
		entry.Type = consulapi.ACLManagementType
	} else {
		entry.Type = consulapi.ACLClientType
	}

	if auo.Raw != "" {
		rules, err := a.ReadRawAcl(auo.Raw)
		if err != nil {
			return err
		}
		entry.Rules = rules
	} else {
		rules, err := a.GetRulesString(auo.ConfigRules)
		if err != nil {
			return err
		}

		entry.Rules = rules
	}

	writeOpts := a.WriteOptions()
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
