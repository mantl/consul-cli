package action

import (
	"flag"
	"fmt"
)

type healthService struct {
	tag     string
	passing bool

	*config
}

func HealthServiceAction() Action {
	return &healthService{
		config: &gConfig,
	}
}

func (h *healthService) CommandFlags() *flag.FlagSet {
	f := h.newFlagSet(FLAG_DATACENTER, FLAG_OUTPUT, FLAG_CONSISTENCY)

	f.StringVar(&h.tag, "tag", "", "Service tag to filter on")
	f.BoolVar(&h.passing, "passing", false, "Only return passing checks")

	return f
}

func (h *healthService) Run(args []string) error {
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

	s, _, err := client.Service(service, h.tag, h.passing, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(s)
}
