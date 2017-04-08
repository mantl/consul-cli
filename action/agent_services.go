package action

import (
	"flag"
)

type agentServices struct {
	*config
}

func AgentServicesAction() Action {
	return &agentServices{
		config: &gConfig,
	}
}

func (a *agentServices) CommandFlags() *flag.FlagSet {
	return a.newFlagSet(FLAG_NONE)
}

func (a *agentServices) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	config, err := client.Services()
	if err != nil {
		return err
	}

	return a.Output(config)
}
