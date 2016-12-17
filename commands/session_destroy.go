package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Destroy functions

func newSessionDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy <sessionId>",
		Short: "Destroy a session",
		Long:  "Destroy a session",
		RunE:  sessionDestroy,
	}

	addDatacenterOption(cmd)

	return cmd
}

func sessionDestroy(cmd *cobra.Command, args []string) error {
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

	_, err = client.Destroy(sessionid, writeOpts)
	if err != nil {
		return err
	}

	return nil
}

