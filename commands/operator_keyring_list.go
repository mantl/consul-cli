package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// keyring list functions

func newOperatorKeyringListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List gossip keys installed",
		Long:  "List gossip keys installed",
		RunE:  operatorKeyringList,
	}

	return cmd
}

func operatorKeyringList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	r, err := client.KeyringList(queryOpts)
	if err != nil {
		return err
	}

	return output(r)
}
