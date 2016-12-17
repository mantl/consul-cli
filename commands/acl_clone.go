package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Clone functions

func newAclCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <token>",
		Short: "Create a new token from an existing one",
		Long:  "Create a new token from an existing one",
		RunE:  aclClone,
	}

	return cmd
}

func aclClone(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	newid, _, err := client.Clone(args[0], writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(newid)

	return nil
}
