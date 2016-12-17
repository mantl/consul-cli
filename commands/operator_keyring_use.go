package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// keyring use functions

func newOperatorKeyringUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <key>",
		Short: "Change the primary gossip encryption key",
		Long:  "Change the primary gossip encryption key",
		RunE:  operatorKeyringUse,
	}

	return cmd
}

func operatorKeyringUse(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Encryption key must be specified")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	return client.KeyringUse(args[0], writeOpts)
}
