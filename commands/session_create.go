package commands

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type SessionCreateOptions struct {
	LockDelay time.Duration
	NodeName  string
	Name      string
	Checks    []string
	Behavior  string
	Ttl       time.Duration
}

func (s *Session) AddCreateSub(cmd *cobra.Command) {
	sco := &SessionCreateOptions{}

	createCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new session",
		Long:  "Create a new session",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Create(args, sco)
		},
	}

	oldCreateCmd := &cobra.Command{
		Use:        "session-create",
		Short:      "Create a new session",
		Long:       "Create a new session",
		Deprecated: "Use session create",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Create(args, sco)
		},
	}

	createCmd.Flags().DurationVar(&sco.LockDelay, "lock-delay", 0, "Lock delay as a duration string")
	createCmd.Flags().StringVar(&sco.Name, "name", "", "Session name")
	createCmd.Flags().StringVar(&sco.NodeName, "node", "", "Node to register session")
	createCmd.Flags().Var((funcVar)(func(s string) error {
		if sco.Checks == nil {
			sco.Checks = make([]string, 0, 1)
		}

		sco.Checks = append(sco.Checks, s)
		return nil
	}), "checks", "Check to associate with session. Can be mulitple")
	createCmd.Flags().StringVar(&sco.Behavior, "behavior", "release", "Lock behavior when session is invalidated. One of release or delete")
	createCmd.Flags().DurationVar(&sco.Ttl, "ttl", 15*time.Second, "Session Time To Live as a duration string")
	s.AddDatacenterOption(createCmd)

	oldCreateCmd.Flags().DurationVar(&sco.LockDelay, "lock-delay", 0, "Lock delay as a duration string")
	oldCreateCmd.Flags().StringVar(&sco.Name, "name", "", "Session name")
	oldCreateCmd.Flags().StringVar(&sco.NodeName, "node", "", "Node to register session")
	oldCreateCmd.Flags().Var((funcVar)(func(s string) error {
		if sco.Checks == nil {
			sco.Checks = make([]string, 0, 1)
		}

		sco.Checks = append(sco.Checks, s)
		return nil
	}), "checks", "Check to associate with session. Can be mulitple")
	oldCreateCmd.Flags().StringVar(&sco.Behavior, "behavior", "release", "Lock behavior when session is invalidated. One of release or delete")
	oldCreateCmd.Flags().DurationVar(&sco.Ttl, "ttl", 15*time.Second, "Session Time To Live as a duration string")
	s.AddDatacenterOption(oldCreateCmd)

	cmd.AddCommand(createCmd)

	s.AddCommand(oldCreateCmd)
}

func (s *Session) Create(args []string, sco *SessionCreateOptions) error {
	// Work around Consul API bug that drops LockDelay == 0
	if sco.LockDelay == 0 {
		sco.LockDelay = time.Nanosecond
	}

	client, err := s.Session()
	if err != nil {
		return err
	}

	writeOpts := s.WriteOptions()

	se := &consulapi.SessionEntry{
		Name:      sco.Name,
		Node:      sco.NodeName,
		Checks:    sco.Checks,
		LockDelay: sco.LockDelay,
		Behavior:  sco.Behavior,
		TTL:       sco.Ttl.String(),
	}

	session, _, err := client.Create(se, writeOpts)
	if err != nil {
		return err
	}

	fmt.Fprintln(s.Out, session)

	return nil
}
