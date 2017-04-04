package action

import (
	"flag"
)

type statusPeers struct {
	*config
}

func StatusPeersAction() Action {
	return &statusPeers{
		config: &gConfig,
	}
}

func (s *statusPeers) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addOutputFlags(f, false)

	return f
}

func (s *statusPeers) Run(args []string) error {
	client, err := s.newStatus()
	if err != nil {
		return err
	}

	l, err := client.Peers()
	if err != nil {
		return err
	}

	return s.Output(l)
}
