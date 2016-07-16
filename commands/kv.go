package commands

import (
	"github.com/spf13/cobra"
)

type Kv struct {
	*Cmd
}

func (root *Cmd) initKv() {
	k := Kv{Cmd: root}

	kvCmd := &cobra.Command{
		Use:   "kv",
		Short: "Consul /kv endpoint interface",
		Long:  "Consul /kv endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	k.AddBulkloadSub(kvCmd)
	k.AddDeleteSub(kvCmd)
	k.AddKeysSub(kvCmd)
	k.AddLockSub(kvCmd)
	k.AddReadSub(kvCmd)
	k.AddUnlockSub(kvCmd)
	k.AddWatchSub(kvCmd)
	k.AddWriteSub(kvCmd)

	k.AddCommand(kvCmd)
}
