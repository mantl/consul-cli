package action

import (
	"flag"
)

type agentMembers struct {
	wan bool
	*config
}

func AgentMembersAction() Action {
	return &agentMembers{
		config: &gConfig,
	}
}

func (a *agentMembers) CommandFlags() *flag.FlagSet {
	f := a.newFlagSet(FLAG_OUTPUT)

	f.BoolVar(&a.wan, "wan", false, "Get list of WAN members instead of LAN")

	return f
}

func (a *agentMembers) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	ms, err := client.Members(a.wan)
	if err != nil {
		return err
	}

	return a.Output(ms)
}
