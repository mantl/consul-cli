package commands

import (
	"github.com/spf13/cobra"
)

func (a *Acl) AddListSub(c *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all active ACL tokens",
		Long:  "List all active ACL tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.List(args)
		},
	}

	oldListCmd := &cobra.Command{
		Use:        "acl-list",
		Short:      "List all active ACL tokens",
		Long:       "List all active ACL tokens",
		Deprecated: "Use acl list",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.List(args)
		},
	}

	a.AddTemplateOption(listCmd)
	c.AddCommand(listCmd)

	a.AddCommand(oldListCmd)
}

func (a *Acl) List(args []string) error {
	client, err := a.ACL()
	if err != nil {
		return err
	}

	queryOpts := a.QueryOptions()
	acls, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return a.Output(acls)
}
