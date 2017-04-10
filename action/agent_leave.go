package action

import (
	"flag"
)

type agentLeave struct {
	*config
}

func AgentLeaveAction() Action {
	return &agentLeave{
		config: &gConfig,
	}
}

func (a *agentLeave) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_NONE)
}

func (a *agentLeave) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	return client.Leave()
}
