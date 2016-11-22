package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (a *Acl) AddDestroySub(c *cobra.Command) {
	destroyCmd := &cobra.Command{
		Use:   "destroy <token>",
		Short: "Destroy an ACL",
		Long:  "Destroy an ACL",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Destroy(args)
		},
	}

	oldDestroyCmd := &cobra.Command{
		Use:        "acl-destroy <token>",
		Short:      "Destroy an ACL",
		Long:       "Destroy an ACL",
		Deprecated: "Use acl destroy",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Destroy(args)
		},
	}

	c.AddCommand(destroyCmd)

	a.AddCommand(oldDestroyCmd)
}

func (a *Acl) Destroy(args []string) error {
	if !a.CheckIdArg(args) {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	client, err := a.ACL()
	if err != nil {
		return err
	}

	writeOpts := a.WriteOptions()
	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
