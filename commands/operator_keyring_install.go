package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// keyring install functions

func newOperatorKeyringInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <key> [<key>]",
		Short: "Install a new gossip key into the cluster",
		Long:  "Install a new gossip key into the cluster",
		RunE:  operatorKeyringInstall,
	}

	return cmd
}

func operatorKeyringInstall(cmd *cobra.Command, args []string) error {
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
		if err := client.KeyringInstall(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
