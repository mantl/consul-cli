package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type checkRegister struct {
	serviceId string

	*config
}

func CheckRegisterAction() Action {
	return &checkRegister{
		config: &gConfig,
	}
}

func (c *checkRegister) CommandFlags() *flag.FlagSet {
	f := c.newFlagSet(FLAG_RAW)

	c.addCheckFlags(f)

	f.StringVar(&c.serviceId, "service-id", "", "Service ID to associate check")

	return f
}

func (c *checkRegister) Run(args []string) error {
	var check consulapi.AgentCheckRegistration

	if c.raw.isSet() {
		if err := c.raw.readJSON(&check); err != nil {
			return err
		}
	} else {
		if len(args) != 1 {
			return fmt.Errorf("A single check id must be specified")
		}
		checkName := args[0]

		checkCount := 0
		if c.check.http != "" {
			checkCount = checkCount + 1
		}
		if c.check.script != "" {
			checkCount = checkCount + 1
		}
		if c.check.ttl != "" {
			checkCount = checkCount + 1
		}
		if c.check.tcp != "" {
			checkCount = checkCount + 1
		}

		if checkCount > 1 {
			return fmt.Errorf("Only one of --http, --script, --tcp or --ttl can be specified")
		}

		check = consulapi.AgentCheckRegistration{
			ID:        c.check.id,
			Name:      checkName,
			ServiceID: c.serviceId,
			Notes:     c.check.notes,
			AgentServiceCheck: consulapi.AgentServiceCheck{
				Script:            c.check.script,
				HTTP:              c.check.http,
				TCP:               c.check.tcp,
				Interval:          c.check.interval,
				TTL:               c.check.ttl,
				TLSSkipVerify:     c.check.skipVerify,
				DockerContainerID: c.check.dockerId,
				Shell:             c.check.shell,
				Notes:             c.check.notes,
				DeregisterCriticalServiceAfter: c.check.deregCrit,
			},
		}
	}

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.CheckRegister(&check)
	if err != nil {
		return err
	}

	return nil
}
