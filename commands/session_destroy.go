package commands

import (

	"github.com/spf13/cobra"
)


func (s *Session) AddDestroySub(cmd *cobra.Command) {
	destroyCmd := &cobra.Command{
		Use: "destroy <sessionId>",
		Short: "Destroy a session",
		Long: "Destroy a session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Destroy(args)
		},
	}

	oldDestroyCmd := &cobra.Command{
		Use: "session-destroy <sessionId>",
		Short: "Destroy a session",
		Long: "Destroy a session",
		Deprecated: "Use session destroy",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Destroy(args)
		},
	}

	s.AddDatacenterOption(destroyCmd)
	s.AddDatacenterOption(oldDestroyCmd)

	cmd.AddCommand(destroyCmd)

	s.AddCommand(oldDestroyCmd)
}

func (s *Session) Destroy(args []string) error {
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

	_, err = sessionClient.Destroy(sessionid, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
