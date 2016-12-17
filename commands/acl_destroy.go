package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Destroy functions

func newAclDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy <token>",
		Short: "Destroy an ACL",
		Long:  "Destroy an ACL",
		RunE:  aclDestroy,
	}

	return cmd
}

func aclDestroy(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()
	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
