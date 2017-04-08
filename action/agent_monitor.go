package action

import (
	"flag"
	"fmt"
)

type agentMonitor struct {
	loglevel string
	*config
}

func AgentMonitorAction() Action {
	return &agentMonitor{
		config: &gConfig,
	}
}

func (a *agentMonitor) CommandFlags() *flag.FlagSet {
	f := a.newFlagSet(FLAG_NONE)

	f.StringVar(&a.loglevel, "loglevel", "", "Log level to filter on. Default is info")

	return f
}

func (a *agentMonitor) Run(args []string) error {
	client, err := a.newAgent()
	if err != nil {
		return err
	}

	outputChan, err := client.Monitor(
		a.loglevel,
		nil, // XXX - Set up interrupts and stop things gracefully
		nil, // query options. No documentation on what they are used for. Leave nil for now.
	)
	if err != nil {
		return err
	}

	for s := range outputChan {
		fmt.Printf(s)
	}

	return nil
}
