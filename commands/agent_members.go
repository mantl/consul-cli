package commands

import (
	"github.com/spf13/cobra"
)

type AgentMembersOptions struct {
	wanFlag bool
}

func (a *Agent) AddMembersSub(c *cobra.Command) {
	amo := &AgentMembersOptions{}

	membersCmd := &cobra.Command{
		Use:   "members",
		Short: "Get the members as seen by the serf agent",
		Long:  "Get the members as seen by the serf agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Members(args, amo)
		},
	}

	oldMembersCmd := &cobra.Command{
		Use:        "agent-members",
		Short:      "Get the members as seen by the serf agent",
		Long:       "Get the members as seen by the serf agent",
		Deprecated: "Use agent members",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Members(args, amo)
		},
	}

	membersCmd.Flags().BoolVar(&amo.wanFlag, "wan", false, "Get list of WAN members instead of LAN")
	oldMembersCmd.Flags().BoolVar(&amo.wanFlag, "wan", false, "Get list of WAN members instead of LAN")

	a.AddTemplateOption(membersCmd)
	c.AddCommand(membersCmd)

	a.AddCommand(oldMembersCmd)
}

func (a *Agent) Members(args []string, amo *AgentMembersOptions) error {
	client, err := a.Agent()
	if err != nil {
		return err
	}

	ms, err := client.Members(amo.wanFlag)
	if err != nil {
		return err
	}

	return a.Output(ms)
}
