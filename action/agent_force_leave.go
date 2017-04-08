package action

import (
	"flag"
	"fmt"
)

type agentForceLeave struct {
	*config
}

func AgentForceLeaveAction() Action {
	return &agentForceLeave{
		config: &gConfig,
	}
}

func (a *agentForceLeave) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_NONE)
}

func (a *agentForceLeave) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Name not provided")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	client, err := a.newAgent()
	if err != nil {
		return err
	}

	return client.ForceLeave(args[0])
}
