package action

import (
	"flag"
	"fmt"
)

type healthChecks struct {
	*config
}

func HealthChecksAction() Action {
	return &healthChecks{
		config: &gConfig,
	}
}

func (h *healthChecks) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	h.addDatacenterFlag(f)
	h.addOutputFlags(f, false)
	h.addConsistencyFlags(f)

	return f
}

func (h *healthChecks) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	client, err := h.newHealth()
	if err != nil {
		return err
	}

	queryOpts := h.queryOptions()

	checks, _, err := client.Checks(service, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(checks)
}
