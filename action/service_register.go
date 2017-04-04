package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type serviceRegister struct {
	checks []map[string]interface{}

	*config
}

func ServiceRegisterAction() Action {
	return &serviceRegister{
		config: &gConfig,
		checks: []map[string]interface{}{},
	}
}

func (s *serviceRegister) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	s.addServiceFlags(f)
	s.addRawFlag(f)

	msv := newMapSliceValue(&s.checks)

	f.Var(msv, "check", "Begin a new check definition")
	f.Var(newMapValue(msv, "http", "string"), "http", "A URL to GET every interval")
	f.Var(newMapValue(msv, "script", "string"), "script", "A script to run every interval")
	f.Var(newMapValue(msv, "ttl", "string"), "ttl", "Fail if TTL expires before service checks in")
	f.Var(newMapValue(msv, "interval", "string"), "interval", "Interval between checks")
	f.Var(newMapValue(msv, "notes", "string"), "notes", "Description of the check")
	f.Var(newMapValue(msv, "docker-id", "string"), "docker-id", "Docker container ID")
	f.Var(newMapValue(msv, "shell", "string"), "shell", "Shell to use inside docker container")
	f.Var(newMapValue(msv, "dereg", "string"), "deregister-crit", "Deregister critical service after this interval")
	f.Var(newMapValue(msv, "skip-verify", "bool"), "skip-verify", "Skip TLS verification for HTTP checks")

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

		checks, err := s.parseChecks()
		if err != nil {
			return err
		}

		service = consulapi.AgentServiceRegistration{
			ID:                s.service.id,
			Name:              serviceName,
			Tags:              s.service.tags,
			Port:              s.service.port,
			Address:           s.service.address,
			Checks:            checks,
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

func (s *serviceRegister) parseChecks() ([]*consulapi.AgentServiceCheck, error) {
	rval := make([]*consulapi.AgentServiceCheck, len(s.checks))

	for i, cs := range s.checks {
		c := new(consulapi.AgentServiceCheck)

		
		if v, ok := cs["script"]; ok { c.Script = v.(string) }
		if v, ok := cs["http"]; ok { c.HTTP = v.(string) }
		if v, ok := cs["ttl"]; ok { c.TTL = v.(string) }
		if v, ok := cs["interval"]; ok { c.Interval = v.(string) }
		if v, ok := cs["notes"]; ok { c.Notes = v.(string) }
		if v, ok := cs["docker-id"]; ok { c.DockerContainerID = v.(string) }
		if v, ok := cs["shell"]; ok { c.Shell = v.(string) }
		if v, ok := cs["dereg"]; ok { c.DeregisterCriticalServiceAfter = v.(string) }
		if v, ok := cs["skip-verify"]; ok { c.TLSSkipVerify = v.(bool) }

		rval[i] = c
	}

	return rval, nil
}
