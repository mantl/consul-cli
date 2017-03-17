package action

import (
	"flag"
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type sessionCreate struct {
	lockDelay time.Duration
	name string
	node string
	checks []string
	behavior string
	ttl time.Duration

	*config
}

func SessionCreateAction() Action {
	return &sessionCreate{
		config: &gConfig,
	}
}

func (s *sessionCreate) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addDatacenterFlag(f)
	s.addRawFlag(f)

	f.DurationVar(&s.lockDelay, "lock-delay", 0, "Lock delay as a duration string")
	f.StringVar(&s.name, "name", "", "Session name")
	f.StringVar(&s.node, "node", "", "Node to register session")
	f.StringVar(&s.behavior, "behavior", "release", "Lock behavior when session is invalidated. One of release or delete")
	f.DurationVar(&s.ttl, "ttl", 15*time.Second, "Session Time To Live as a duration string")
	f.Var(newStringSliceValue(&s.checks), "checks", "Check to associate with session. Can be mulitple")

	return f
}

func (s *sessionCreate) Run(args []string) error {
	var se consulapi.SessionEntry

	if s.raw.isSet() {
		if err := s.raw.readJSON(&se); err != nil {
			return err
		}
	} else {

		// Work around Consul API bug that drops LockDelay == 0
		if s.lockDelay == 0 {
			s.lockDelay = time.Nanosecond
		}

		se = consulapi.SessionEntry{
			Name:      s.name,
			Node:      s.node,
			Checks:    s.checks,
			LockDelay: s.lockDelay,
			Behavior:  s.behavior,
			TTL:       s.ttl.String(),
		}
	}
	writeOpts := s.writeOptions()
	client, err := s.newSession()
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
