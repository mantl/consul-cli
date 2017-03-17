package action

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type serviceRegister struct {
	*config
}

func ServiceRegisterAction() Action {
	return &serviceRegister{
		config: &gConfig,
	}
}

func (s *serviceRegister) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addServiceFlags(f)
	s.addRawFlag(f)

	return f
}

func (s *serviceRegister) Run(args []string) error {
	var service consulapi.AgentServiceRegistration

	if s.raw.isSet() {
		if err := s.raw.readJSON(&service); err != nil {
			return err
		}
	} else {
		switch {
		case len(args) == 0:
			return fmt.Errorf("Service name must be specified")
		case len(args) > 1:
			return fmt.Errorf("Only one service name allowed")
		}
		serviceName := args[0]

		//checkStrings := viper.GetStringSlice("check")
		//checks, err := parseChecks(checkStrings)
		//if err != nil {
		//	return err
		//}

		service = consulapi.AgentServiceRegistration{
			ID:                s.service.id,
			Name:              serviceName,
			Tags:              s.service.tags,
			Port:              s.service.port,
			Address:           s.service.address,
//			Checks:            checks,
			EnableTagOverride: s.service.overrideTag,
		}
	}

	agent, err := s.newAgent()
	if err != nil {
		return err
	}

	if err := agent.ServiceRegister(&service); err != nil {
		return err
	}

	return nil
}

func parseChecks(checks []string) ([]*consulapi.AgentServiceCheck, error) {
	rval := make([]*consulapi.AgentServiceCheck, len(checks))

	for i, s := range checks {
		parts := strings.Split(s, ":")
		numParts := len(parts)

		switch parts[0] {
		case "http":
			// Order is http:interval:url:[tlsskip]
			switch {
			case numParts < 3 || numParts > 4:
				return nil, fmt.Errorf("HTTP check format is http:interval:url:[tlsskipverify]")
			case numParts == 3:
				rval[i] = &consulapi.AgentServiceCheck{
					HTTP:     parts[2],
					Interval: parts[1],
				}
			case numParts == 4:
				rval[i] = &consulapi.AgentServiceCheck{
					HTTP:     parts[2],
					Interval: parts[1],
					//TLSSkipVerify: parts[3],
				}
			}
		case "script":
			// Order is script:interval:command
			if numParts != 3 {
				return nil, fmt.Errorf("Script check format is script:interval:command_path")
			}

			rval[i] = &consulapi.AgentServiceCheck{
				Script:   parts[2],
				Interval: parts[1],
			}
		case "ttl":
			// Order is ttl:interval
			if numParts != 2 {
				return nil, fmt.Errorf("TTL check format is ttl:interval")
			}

			rval[i] = &consulapi.AgentServiceCheck{
				TTL: parts[1],
			}
		}
	}

	return rval, nil
}
func parseCheckConfig(s string) (*consulapi.AgentServiceCheck, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, errors.New("Cannot specify empty check")
	}

	var checkType, checkInterval, checkString string
	parts := strings.Split(s, ":")
	partLen := len(parts)
	switch {
	case partLen == 2:
		checkType, checkInterval = parts[0], parts[1]
		checkString = ""
	case partLen >= 3:
		checkType, checkInterval, checkString = parts[0], parts[1], strings.Join(parts[2:], ":")
	default:
		return nil, fmt.Errorf("Invalid check definition '%s'", s)
	}

	switch strings.ToLower(checkType) {
	case "http":
		if checkString == "" {
			return nil, errors.New("Must provide URL with HTTP check")
		}

		return &consulapi.AgentServiceCheck{
			HTTP:     checkString,
			Interval: checkInterval,
		}, nil
	case "script":
		if checkString == "" {
			return nil, errors.New("Must provide command path with script check")
		}

		return &consulapi.AgentServiceCheck{
			Script:   checkString,
			Interval: checkInterval,
		}, nil
	case "ttl":
		return &consulapi.AgentServiceCheck{
			TTL: checkInterval,
		}, nil
	}

	return nil, fmt.Errorf("Invalid check type '%s'", checkType)
}
