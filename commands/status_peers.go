package commands

import (
	"github.com/spf13/cobra"
)

func (s *Status) AddPeersSub(cmd *cobra.Command) {
	peersCmd := &cobra.Command{
		Use:   "peers",
		Short: "Get the current Raft peers",
		Long:  "Get the current Raft peers",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Peers(args)
		},
	}

	oldPeersCmd := &cobra.Command{
		Use:        "status-peers",
		Short:      "Get the current Raft peers",
		Long:       "Get the current Raft peers",
		Deprecated: "Use status peers",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Peers(args)
		},
	}

	s.AddTemplateOption(peersCmd)
	cmd.AddCommand(peersCmd)

	s.AddCommand(oldPeersCmd)
}

func (s *Status) Peers(args []string) error {
	client, err := s.Status()
	if err != nil {
		return err
	}


	l, err := client.Peers()
	if err != nil {
		return err
	}

	return s.Output(l)
}
