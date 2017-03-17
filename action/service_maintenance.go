package action

import (
	"flag"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type serviceMaintenance struct {
	enabled bool
	reason string

	*config
}

func ServiceMaintenanceAction() Action {
	return &serviceMaintenance{
		config: &gConfig,
	}
}

func (s *serviceMaintenance) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.BoolVar(&s.enabled, "enabled", true, "Boolean value for maintenance mode")
	f.StringVar(&s.reason, "reason", "", "Reason for entering maintenance mode")

	return f
}

func (s *serviceMaintenance) Run(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	agent, err := s.newAgent()
	if err != nil {
		return err
	}

	var result error

	for _, id := range args {
		if s.enabled {
			if err := agent.EnableServiceMaintenance(id, s.reason); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			if err := agent.DisableServiceMaintenance(id); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return result
}
