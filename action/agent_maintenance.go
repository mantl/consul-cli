package action

import (
	"flag"
)

type agentMaintenance struct {
	enabled bool
	reason  string
	*config
}

func AgentMaintenanceAction() Action {
	return &agentMaintenance{
		config: &gConfig,
	}
}

func (a *agentMaintenance) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&a.enabled, "enabled", true, "Boolean value for maintenance mode")
	f.StringVar(&a.reason, "reason", "", "Reason for entering maintenance mode")

	return f
}

func (a *agentMaintenance) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	if a.enabled {
		return client.EnableNodeMaintenance(a.reason)
	}

	return client.DisableNodeMaintenance()
}
