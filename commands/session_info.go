package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Info functions

func newSessionInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <sessionId>",
		Short: "Get information on a session",
		Long:  "Get information on a session",
		RunE:  sessionInfo,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func sessionInfo(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	sessionid := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	session, _, err := client.Info(sessionid, queryOpts)
	if err != nil {
		return err
	}

	return output(session)
}
