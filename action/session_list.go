package action

import (
	"flag"
)

type sessionList struct {
	*config
}

func SessionListAction() Action {
	return &sessionList{
		config: &gConfig,
	}
}

func (s *sessionList) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_DATACENTER, FLAG_OUTPUT, FLAG_CONSISTENCY)
}

func (s *sessionList) Run(args []string) error {
	client, err := s.newSession()
	if err != nil {
		return err
	}

	queryOpts := s.queryOptions()

	sessions, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return s.Output(sessions)
}
