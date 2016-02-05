package commands

import (
	"github.com/spf13/cobra"
)

func (s *Session) AddListSub(cmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List active sessions for a datacenter",
		Long:  "List active sessions for a datacenter",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.List(args)
		},
	}

	oldListCmd := &cobra.Command{
		Use:        "session-list",
		Short:      "List active sessions for a datacenter",
		Long:       "List active sessions for a datacenter",
		Deprecated: "Use session list",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.List(args)
		},
	}

	s.AddDatacenterOption(listCmd)
	s.AddTemplateOption(listCmd)
	s.AddDatacenterOption(oldListCmd)

	cmd.AddCommand(listCmd)

	s.AddCommand(oldListCmd)
}

func (s *Session) List(args []string) error {
	client, err := s.Session()
	if err != nil {
		return err
	}

	queryOpts := s.QueryOptions()

	sessions, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return s.Output(sessions)
}
