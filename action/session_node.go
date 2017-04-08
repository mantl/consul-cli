package action

import (
	"flag"
	"fmt"
)

type sessionNode struct {
	*config
}

func SessionNodeAction() Action {
	return &sessionNode{
		config: &gConfig,
	}
}

func (s *sessionNode) CommandFlags() *flag.FlagSet {
	return s.newFlagSet(FLAG_DATACENTER, FLAG_OUTPUT, FLAG_CONSISTENCY)
}

func (s *sessionNode) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single session Id must be specified")
	}
	node := args[0]

	client, err := s.newSession()
	if err != nil {
		return err
	}

	queryOpts := s.queryOptions()

	sessions, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return s.Output(sessions)
}
