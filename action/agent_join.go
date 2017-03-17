package action

import (
	"flag"
	"fmt"
)

type agentJoin struct {
	wan bool
	*config
}

func AgentJoinAction() Action {
	return &agentJoin{
		config: &gConfig,
	}
}

func (a *agentJoin) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&a.wan, "wan", false, "Get list of WAN join instead of LAN")

	return f
}
func (a *agentJoin) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one name allowed")
	}

	client, err := a.newAgent()
	if err != nil {
		return err
	}

	return client.Join(args[0], a.wan)
}
