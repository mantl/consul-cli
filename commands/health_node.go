package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Node functions

func newHealthNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get the health info for a node",
		Long:  "Get the health info for a node",
		RunE:  healthNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthNode(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node name allowed")
	}
	node := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	n, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return output(n)
}
