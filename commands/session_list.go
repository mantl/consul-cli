package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// List functions

func newSessionListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List active sessions for a datacenter",
		Long:  "List active sessions for a datacenter",
		RunE:  sessionList,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func sessionList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	sessions, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return output(sessions)
}
