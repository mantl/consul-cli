package commands

import (
	"github.com/spf13/cobra"
)

func (s *Status) AddLeaderSub(cmd *cobra.Command) {
	leaderCmd := &cobra.Command{
		Use: "leader",
		Short: "Get the current Raft leader",
		Long: "Get the current Raft leader",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Leader(args)
		},
	}

	oldLeaderCmd := &cobra.Command{
		Use: "status-leader",
		Short: "Get the current Raft leader",
		Long: "Get the current Raft leader",
		Deprecated: "Use status leader",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Leader(args)
		},
	}

	s.AddTemplateOption(leaderCmd)
	cmd.AddCommand(leaderCmd)

	s.AddCommand(oldLeaderCmd)
}

func (s *Status) Leader(args []string) error {
	client, err := s.Client()
	if err != nil {
		return err
	}

	statusClient := client.Status()

	l, err := statusClient.Leader()
	if err != nil {
		return err
	}

	return s.Output(l)
}
