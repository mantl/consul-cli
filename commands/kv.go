package commands

import (
	"github.com/spf13/cobra"
)

func newKvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kv",
		Short: "Consul /kv endpoint interface",
		Long:  "Consul /kv endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newKvBulkloadCommand())
	cmd.AddCommand(newKvDeleteCommand())
	cmd.AddCommand(newKvKeysCommand())
	cmd.AddCommand(newKvLockCommand())
	cmd.AddCommand(newKvReadCommand())
	cmd.AddCommand(newKvUnlockCommand())
	cmd.AddCommand(newKvWatchCommand())
	cmd.AddCommand(newKvWriteCommand())

	return cmd
}

