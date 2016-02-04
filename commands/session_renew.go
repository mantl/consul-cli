package commands

import (
	"github.com/spf13/cobra"
)

func (s *Session) AddRenewSub(cmd *cobra.Command) {
	renewCmd := &cobra.Command{
		Use: "renew <sessionId>",
		Short: "Renew the given session",
		Long: "Renew the given session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Renew(args)
		},
	}

	oldRenewCmd := &cobra.Command{
		Use: "session-renew <sessionId>",
		Short: "Renew the given session",
		Long: "Renew the given session",
		Deprecated: "Use session renew",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Renew(args)
		},
	}

	s.AddDatacenterOption(renewCmd)
	s.AddTemplateOption(renewCmd)
	s.AddDatacenterOption(oldRenewCmd)

	cmd.AddCommand(renewCmd)

	s.AddCommand(oldRenewCmd)
}

func (s *Session) Renew(args []string) error {
	if err := s.CheckIdArg(args); err != nil {
		return err
	}
	sessionid := args[0]

	client, err := s.Client()
	if err != nil {
		return err
	}

	writeOpts := s.WriteOptions()
	sessionClient := client.Session()

	session, _, err := sessionClient.Renew(sessionid, writeOpts)
	if err != nil {
		return err
	}

	if session != nil {
		s.Output(session)
	}

	return nil
}
