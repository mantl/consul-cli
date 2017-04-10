package action

import (
	"flag"
	"fmt"
)

type sessionRenew struct {
	*config
}

func SessionRenewAction() Action {
	return &sessionRenew{
		config: &gConfig,
	}
}

func (s *sessionRenew) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_DATACENTER, FLAG_OUTPUT)
}

func (s *sessionRenew) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	sessionid := args[0]

	client, err := s.newSession()
	if err != nil {
		return err
	}

	writeOpts := s.writeOptions()

	session, _, err := client.Renew(sessionid, writeOpts)
	if err != nil {
		return err
	}

	if session != nil {
		s.Output(session)
	}

	return nil
}
