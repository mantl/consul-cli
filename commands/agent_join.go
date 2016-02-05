package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type AgentJoinOptions struct {
	wanFlag bool
}

func (a *Agent) AddJoinSub(c *cobra.Command) {
	ajo := &AgentJoinOptions{}

	joinCmd := &cobra.Command{
		Use:   "join",
		Short: "Trigger the local agent to join a node",
		Long:  "Trigger the local agent to join a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Join(args, ajo)
		},
	}

	oldJoinCmd := &cobra.Command{
		Use:        "agent-join",
		Short:      "Trigger the local agent to join a node",
		Long:       "Trigger the local agent to join a node",
		Deprecated: "Use agent join",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Join(args, ajo)
		},
	}

	joinCmd.Flags().BoolVar(&ajo.wanFlag, "wan", false, "Get list of WAN join instead of LAN")
	oldJoinCmd.Flags().BoolVar(&ajo.wanFlag, "wan", false, "Get list of WAN join instead of LAN")

	a.AddTemplateOption(joinCmd)
	c.AddCommand(joinCmd)

	a.AddCommand(oldJoinCmd)
}

func (a *Agent) Join(args []string, ajo *AgentJoinOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one name allowed")
	}

	client, err := a.Agent()
	if err != nil {
		return err
	}

	return client.Join(args[0], ajo.wanFlag)
}
