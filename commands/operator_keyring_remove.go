package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// keyring remove functions

func newOperatorKeyringRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <key> [<key>]",
		Short: "Remove gossip keys from the cluster",
		Long:  "Remove gossip keys from the cluster",
		RunE:  operatorKeyringRemove,
	}

	return cmd
}

func operatorKeyringRemove(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	if len(args) < 1 {
		return fmt.Errorf("At least one gossip key must be specified")
	}

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	var result error

	for _, k := range args {
		if err := client.KeyringRemove(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

