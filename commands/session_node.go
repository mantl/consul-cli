package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Node functions

func newSessionNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get active sessions for a node",
		Long:  "Get active sessions for a node",
		RunE:  sessionNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func sessionNode(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	node := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	sessions, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return output(sessions)
}
