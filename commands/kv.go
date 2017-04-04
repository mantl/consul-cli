package commands

import (
	"github.com/spf13/cobra"

	"github.com/ChrisAubuchon/consul-cli/action"
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

func newKvBulkloadCommand() *cobra.Command {
	k := action.KvBulkloadAction()

        cmd := &cobra.Command{
                Use:   "bulkload",
                Short: "Bulkload value to the K/V store",
                Long:  "Bulkload value to the K/V store",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvDeleteCommand() *cobra.Command {
	k := action.KvDeleteAction()

        cmd := &cobra.Command{
                Use:   "delete <path>",
                Short: "Delete a given path from the K/V",
                Long:  "Delete a given path from the K/V",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvKeysCommand() *cobra.Command {
	k := action.KvKeysAction()

        cmd := &cobra.Command{
                Use:   "keys <path>",
                Short: "List K/V keys",
                Long:  "List K/V keys",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvLockCommand() *cobra.Command {
	k := action.KvLockAction()

        cmd := &cobra.Command{
                Use:   "lock <path>",
                Short: "Acquire a lock on a given path",
                Long:  "Acquire a lock on a given path",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvReadCommand() *cobra.Command {
	k := action.KvReadAction()

        cmd := &cobra.Command{
                Use:   "read <path>",
                Short: "Read a value from a given path",
                Long:  "Read a value from a given path",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvUnlockCommand() *cobra.Command {
	k := action.KvUnlockAction()

        cmd := &cobra.Command{
                Use:   "unlock <path>",
                Short: "Release a lock on a given path",
                Long:  "Release a lock on a given path",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvWatchCommand() *cobra.Command {
	k := action.KvWatchAction()

        cmd := &cobra.Command{
                Use:   "watch <path>",
                Short: "Watch for changes to a K/V path",
                Long:  "Watch for changes to a K/V path",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

func newKvWriteCommand() *cobra.Command {
	k := action.KvWriteAction()

        cmd := &cobra.Command{
                Use:   "write <path> <value>",
                Short: "Write a value to a given path",
                Long:  "Write a value to a given path",
		RunE: func (cmd *cobra.Command, args []string) error {
			return k.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(k.CommandFlags())

	return cmd
}

