package commands

import (
	"github.com/spf13/cobra"
)

func (a *Acl) AddInfoSub(c *cobra.Command) {
	infoCmd := &cobra.Command{
		Use: "info <token>",
		Short: "Query information about an ACL token",
		Long: "Query information about an ACL token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Info(args)
		},
	}

	oldInfoCmd := & cobra.Command{
		Use: "acl-info <token>",
		Short: "Query information about an ACL token",
		Long: "Query information about an ACL token",
		Deprecated: "Use acl info",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Info(args)
		},
	}

	a.AddTemplateOption(infoCmd)
	c.AddCommand(infoCmd)

	a.AddCommand(oldInfoCmd)
}


func (a *Acl) Info(args []string) error {
	if err := a.CheckIdArg(args); err != nil {
		return err
	}
	id := args[0]

	consul, err := a.Client()
	if err != nil {
		return err
	}
	client := consul.ACL()

	queryOpts := a.QueryOptions()
	acl, _, err := client.Info(id, queryOpts)
	if err != nil {
		return err
	}

	return a.Output(acl)
}

