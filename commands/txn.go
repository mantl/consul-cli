package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
)

func newTxnCommand() *cobra.Command {
	t := action.TxnAction()

	cmd := &cobra.Command{
		Use:   "txn",
		Short: "Consul /txn endpoint interface",
		Long:  "Consul /txn endpoint interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			return t.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(t.CommandFlags())

	return cmd
}
