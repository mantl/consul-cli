package action

import (
	"flag"
	"fmt"
)

type sessionInfo struct {
	*config
}

func SessionInfoAction() Action {
	return &sessionInfo{
		config: &gConfig,
	}
}

func (s *sessionInfo) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addDatacenterFlag(f)
	s.addOutputFlags(f, false)
	s.addConsistencyFlags(f)

	return f
}

func (s *sessionInfo) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	sessionid := args[0]

	client, err := s.newSession()
	if err != nil {
		return err
	}

	queryOpts := s.queryOptions()

	session, _, err := client.Info(sessionid, queryOpts)
	if err != nil {
		return err
	}

	return s.Output(session)
}
