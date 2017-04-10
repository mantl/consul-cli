package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
)

func newSnapshotCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Consul /snapshot endpoint interface",
		Long:  "Consul /snapshot endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newSnapshotSaveCommand())
	cmd.AddCommand(newSnapshotRestoreCommand())

	return cmd
}

func newSnapshotSaveCommand() *cobra.Command {
	s := action.SnapshotSaveAction()

	cmd := &cobra.Command{
		Use:   "save snapshot_path",
		Short: "Create a new snapshot",
		Long:  "Create a new snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}

func newSnapshotRestoreCommand() *cobra.Command {
	s := action.SnapshotRestoreAction()

	cmd := &cobra.Command{
		Use:   "restore snapshot_path",
		Short: "Restore a snapshot",
		Long:  "Restore a snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(s.CommandFlags())

	return cmd
}
