package command

import (
	"errors"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type ServiceRegisterCommand struct {
	Meta
	id			string
	tags			[]string
	address			string
	port			int
	checks			consulapi.AgentServiceChecks
}

func (c *ServiceRegisterCommand) Help() string {
	helpText := `
Usage: consul-cli service-register [options] serviceName

  Register a new local service

  If --id is not specified, the serviceName is used. There cannot
be duplicate service IDs per agent however.

  If --address is not specified, the IP address of the local agent
is used.

Options:
` + c.ConsulHelp() + 
`  --id				Service Id
				(default: not set)
  --tag				Service tag. Multiple tags allowed
				(default: not set)
  --address			Service address
				(default: not set)
  --port			Service port
				(default: not set)
  --check			Check to create. Multiple checks can be
				created with the service. The format is:
				[http|script|ttl]:<interval>:<command|url>
`

	return strings.TrimSpace(helpText)
}

func (c *ServiceRegisterCommand) Run(args []string) int {
	flags := c.Meta.FlagSet(false)
	flags.StringVar(&c.id, "id", "", "")
	flags.Var((funcVar)(func(s string) error {
		if c.tags == nil {
			c.tags = make([]string, 0, 1)
		}

		c.tags = append(c.tags, s)
		return nil
	}), "tag", "")
	flags.StringVar(&c.address, "address", "", "")
	flags.IntVar(&c.port, "port", 0, "")
	flags.Var((funcVar)(func(s string) error {
		t, err := ParseCheckConfig(s)
		if err != nil {
			return err
		}

		if c.checks == nil {
			c.checks = make(consulapi.AgentServiceChecks, 0, 1)
		}

		c.checks = append(c.checks, t)
		return nil
	}), "check", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Service name must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	serviceName := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	service := &consulapi.AgentServiceRegistration{
		ID:		c.id,
		Name:		serviceName,
		Tags:		c.tags,
		Port:		c.port,
		Address:	c.address,
		Checks:		c.checks,
	}

	client := consul.Agent()
	err = client.ServiceRegister(service)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *ServiceRegisterCommand) Synopsis() string {
	return "Register a new local service"
}

func ParseCheckConfig(s string) (*consulapi.AgentServiceCheck, error) {
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
			HTTP:		checkString,
			Interval:	checkInterval,
		}, nil
	case "script":
		if checkString == "" {
			return nil, errors.New("Must provide command path with script check")
		}

		return &consulapi.AgentServiceCheck{
			Script:		checkString,
			Interval:	checkInterval,
		}, nil
	case "ttl":
		return &consulapi.AgentServiceCheck{
			TTL:		checkInterval,
		}, nil
	}

	return nil, fmt.Errorf("Invalid check type '%s'", checkType)
}
