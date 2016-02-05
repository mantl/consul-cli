package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *Acl) AddCloneSub(c *cobra.Command) {
	cloneCmd := &cobra.Command{
		Use:   "clone <token>",
		Short: "Create a new token from an existing one",
		Long:  "Create a new token from an existing one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Clone(args)
		},
	}

	oldCloneCmd := &cobra.Command{
		Use:        "acl-clone <token>",
		Short:      "Create a new token from an existing one",
		Long:       "Create a new token from an existing one",
		Deprecated: "Use acl clone",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Clone(args)
		},
	}

	a.AddTemplateOption(cloneCmd)
	c.AddCommand(cloneCmd)

	a.AddCommand(oldCloneCmd)
}

func (a *Acl) Clone(args []string) error {
	if err := a.CheckIdArg(args); err != nil {
		return err
	}

	client, err := a.ACL()
	if err != nil {
		return err
	}

	writeOpts := a.WriteOptions()
	newid, _, err := client.Clone(args[0], writeOpts)
	if err != nil {
		return err
	}

	fmt.Fprintln(a.Out, newid)

	return nil
}
