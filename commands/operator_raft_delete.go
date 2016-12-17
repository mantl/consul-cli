package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Peer delete functions

var raftDeleteLongHelp = `
Remove a Consul server from the Raft configuration

An address is required and should be set to IP:port for the server
to remove. The port number is 8300 unless configured otherwise`

func newOperatorRaftDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <address>",
		Short: "Remove a Consul server from the Raft configuration",
		Long:  raftDeleteLongHelp,
		RunE:  operatorRaftDelete,
	}

	addDatacenterOption(cmd)

	return cmd
}

func operatorRaftDelete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single address argument must be specified")
	}
	address := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	return client.RaftRemovePeerByAddress(address, writeOpts)
}
