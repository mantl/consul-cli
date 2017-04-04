package action

import (
	"flag"
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type serviceDeregister struct {
	*config
}

func ServiceDeregisterAction() Action {
	return &serviceDeregister{
		config: &gConfig,
	}
}

func (s *serviceDeregister) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

func (s *serviceDeregister) Run(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	agent, err := s.newAgent()
	if err != nil {
		return err
	}

	var result error

	for _, id := range args {
		if err := agent.ServiceDeregister(id); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}
