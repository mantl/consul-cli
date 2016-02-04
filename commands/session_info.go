package commands

import (
	"github.com/spf13/cobra"
)

func (s *Session) AddInfoSub(cmd *cobra.Command) {
	infoCmd := &cobra.Command{
		Use: "info <sessionId>",
		Short: "Get information on a session",
		Long: "Get information on a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Info(args)
		},
	}

	oldInfoCmd := &cobra.Command{
		Use: "info <sessionId>",
		Short: "Get information on a session",
		Long: "Get information on a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Info(args)
		},
	}

	s.AddDatacenterOption(infoCmd)
	s.AddTemplateOption(infoCmd)
	s.AddDatacenterOption(oldInfoCmd)

	cmd.AddCommand(infoCmd)

	s.AddCommand(oldInfoCmd)
}

func (s *Session) Info(args []string) error {
	if err := s.CheckIdArg(args); err != nil {
		return err
	}
	sessionid := args[0]

	client, err := s.Client()
	if err != nil {
		return err
	}

	queryOpts := s.QueryOptions()
	sessionClient := client.Session()

	session, _, err := sessionClient.Info(sessionid, queryOpts)
	if err != nil {
		return err
	}

	return s.Output(session)
}
