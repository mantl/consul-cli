package action

import (
	"flag"
	"fmt"
)

type sessionDestroy struct {
	*config
}

func SessionDestroyAction() Action {
	return &sessionDestroy{
		config: &gConfig,
	}
}

func (s *sessionDestroy) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addDatacenterFlag(f)

	return f
}

func (s *sessionDestroy) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	sessionid := args[0]

	client, err := s.newSession()
	if err != nil {
		return err
	}

	writeOpts := s.writeOptions()

	_, err = client.Destroy(sessionid, writeOpts)
	if err != nil {
		return err
	}

	return nil
}

