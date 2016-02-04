package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (h *Health) AddNodeSub(cmd *cobra.Command) {
	nodeCmd := &cobra.Command{
		Use: "node <nodeName>",
		Short: "Get the health info for a node",
		Long: "Get the health info for a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Node(args)
		},
	}

	oldNodeCmd := &cobra.Command{
		Use: "health-node <nodeName>",
		Short: "Get the health info for a node",
		Long: "Get the health info for a node",
		Deprecated: "Use health node",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.Node(args)
		},
	}

	h.AddDatacenterOption(nodeCmd)
	h.AddDatacenterOption(oldNodeCmd)
	h.AddTemplateOption(nodeCmd)

	cmd.AddCommand(nodeCmd)

	h.AddCommand(oldNodeCmd)
}

func (h *Health) Node(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node name allowed")
	}
	node := args[0]

	client, err := h.Client()
	if err != nil {
		return err
	}

	queryOpts := h.QueryOptions()
	healthClient := client.Health()

	n, _, err := healthClient.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(n)
}
