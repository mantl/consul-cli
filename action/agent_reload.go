package action

import (
	"flag"
)

type agentReload struct {
	*config
}

func AgentReloadAction() Action {
	return &agentReload{
		config: &gConfig,
	}
}

func (a *agentReload) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (a *agentReload) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	return client.Reload()
}
