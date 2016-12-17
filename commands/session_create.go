package commands

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Create functions

func newSessionCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new session",
		Long:  "Create a new session",
		RunE:  sessionCreate,
	}

	cmd.Flags().Duration("lock-delay", 0, "Lock delay as a duration string")
	cmd.Flags().String("name", "", "Session name")
	cmd.Flags().String("node", "", "Node to register session")
	cmd.Flags().StringSlice("checks", nil, "Check to associate with session. Can be mulitple")
	cmd.Flags().String("behavior", "release", "Lock behavior when session is invalidated. One of release or delete")
	cmd.Flags().Duration("ttl", 15*time.Second, "Session Time To Live as a duration string")

	addDatacenterOption(cmd)
	addRawOption(cmd)

	return cmd
}

func sessionCreate(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var se consulapi.SessionEntry

	if raw := viper.GetString("raw"); raw != "" {
		if err := readRawJSON(raw, &se); err != nil {
			return err
		}
	} else {

		// Work around Consul API bug that drops LockDelay == 0
		if viper.GetDuration("lock-delay") == 0 {
			viper.Set("lock-delay", time.Nanosecond)
		}

		se = consulapi.SessionEntry{
			Name:      viper.GetString("name"),
			Node:      viper.GetString("node"),
			Checks:    viper.GetStringSlice("checks"),
			LockDelay: viper.GetDuration("lock-delay"),
			Behavior:  viper.GetString("behavior"),
			TTL:       viper.GetDuration("ttl").String(),
		}
	}
	writeOpts := writeOptions()
	client, err := newSession()
	if err != nil {
		return err
	}

	session, _, err := client.Create(&se, writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(session)

	return nil
}
