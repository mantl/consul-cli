package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (s *Session) AddNodeSub(cmd *cobra.Command) {
	nodeCmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get active sessions for a node",
		Long:  "Get active sessions for a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Node(args)
		},
	}

	oldNodeCmd := &cobra.Command{
		Use:        "session-node <nodeName>",
		Short:      "Get active sessions for a node",
		Long:       "Get active sessions for a node",
		Deprecated: "Use session node",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Node(args)
		},
	}

	s.AddDatacenterOption(nodeCmd)
	s.AddTemplateOption(nodeCmd)
	s.AddDatacenterOption(oldNodeCmd)

	cmd.AddCommand(nodeCmd)

	s.AddCommand(oldNodeCmd)
}

func (s *Session) Node(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node name allowed")
	}
	node := args[0]

	client, err := s.Session()
	if err != nil {
		return err
	}

	queryOpts := s.QueryOptions()

	sessions, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return s.Output(sessions)
}
