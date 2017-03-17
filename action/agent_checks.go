package action

import (
	"flag"
)

type agentChecks struct {
	*config
}

func AgentChecksAction() Action {
	return &agentChecks{
		config: &gConfig,
	}
}

func (a *agentChecks) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (a *agentChecks) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	config, err := client.Checks()
	if err != nil {
		return err
	}

	return a.Output(config)
}
