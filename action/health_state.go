package action

import (
	"flag"
	"fmt"
	"strings"
)

type healthState struct {
	*config
}

func HealthStateAction() Action {
	return &healthState{
		config: &gConfig,
	}
}

func (h *healthState) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	h.addDatacenterFlag(f)
	h.addOutputFlags(f, false)
	h.addConsistencyFlags(f)

	return f
}

func (h *healthState) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Check state must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one check state allowed")
	}
	state := strings.ToLower(args[0])

	client, err := h.newHealth()
	if err != nil {
		return err
	}

	queryOpts := h.queryOptions()

	s, _, err := client.State(state, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(s)
}
