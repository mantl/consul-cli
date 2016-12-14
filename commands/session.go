package commands

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newSessionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "session",
		Short: "Consul /session endpoint interface",
		Long:  "Consul /session endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newSessionCreateCommand())
	cmd.AddCommand(newSessionDestroyCommand())
	cmd.AddCommand(newSessionInfoCommand())
	cmd.AddCommand(newSessionListCommand())
	cmd.AddCommand(newSessionNodeCommand())
	cmd.AddCommand(newSessionRenewCommand())

	return cmd
}

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

// Node functions

func newSessionNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get active sessions for a node",
		Long:  "Get active sessions for a node",
		RunE:  sessionNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func sessionNode(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	node := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	sessions, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return output(sessions)
}

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
