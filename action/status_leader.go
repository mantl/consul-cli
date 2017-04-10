package action

import (
	"flag"
)

type statusLeader struct {
	*config
}

func StatusLeaderAction() Action {
	return &statusLeader{
		config: &gConfig,
	}
}

func (s *statusLeader) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_OUTPUT)
}

func (s *statusLeader) Run(args []string) error {
	client, err := s.newStatus()
	if err != nil {
		return err
	}

	l, err := client.Leader()
	if err != nil {
		return err
	}

	return s.Output(l)
}
