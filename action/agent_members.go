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
	f := newFlagSet()

	f.BoolVar(&a.wan, "wan", false, "Get list of WAN members instead of LAN")
	a.addOutputFlags(f, false)

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
