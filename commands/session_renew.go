package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Renew functions

func newSessionRenewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "renew <sessionId>",
		Short: "Renew the given session",
		Long:  "Renew the given session",
		RunE:  sessionRenew,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func sessionRenew(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	sessionid := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newSession()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	session, _, err := client.Renew(sessionid, writeOpts)
	if err != nil {
		return err
	}

	if session != nil {
		output(session)
	}

	return nil
}
