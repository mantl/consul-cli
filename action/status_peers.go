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
	return s.newFlagSet(FLAG_OUTPUT)
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
