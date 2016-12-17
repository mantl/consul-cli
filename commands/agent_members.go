package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Members functions

func newAgentMembersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Short: "Get the members as seen by the serf agent",
		Long:  "Get the members as seen by the serf agent",
		RunE:  agentMembers,
	}

	cmd.Flags().Bool("wan", false, "Get list of WAN members instead of LAN")

	addTemplateOption(cmd)

	return cmd
}

func agentMembers(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	ms, err := client.Members(viper.GetBool("wan"))
	if err != nil {
		return err
	}

	return output(ms)
}
