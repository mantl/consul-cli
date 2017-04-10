package action

import (
	"flag"
)

type agentSelf struct {
	*config
}

func AgentSelfAction() Action {
	return &agentSelf{
		config: &gConfig,
	}
}

func (a *agentSelf) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_OUTPUT)
}

func (a *agentSelf) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	config, err := client.Self()
	if err != nil {
		return err
	}

	return a.Output(config)
}
