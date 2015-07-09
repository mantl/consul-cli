package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type ServiceRegisterCommand struct {
	Meta
	id			string
	tags			[]string
	address			string
	port			int
	checkName		string
	checkScript		string
	checkHttp		string
	checkTtl		string
	checkInterval		string
}

func (c *ServiceRegisterCommand) Help() string {
	helpText := `
Usage: consul-cli service-register [options] serviceName

  Register a new local service

  If --id is not specified, the serviceName is used. There cannot
be duplicate service IDs per agent however.

  If --address is not specified, the IP address of the local agent
is used.

  Only one of --check-http, --check-script and --check-ttl can be
specified.

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
  --check-http			A URL to GET every interval
				(default: not set)
  --check-script		A script to run every interval
				(default: not set)
  --check-ttl			Fail if TTL expires before service 
				checks in
				(default: not set)
  --check-interval		Interval between checks
				(default: not set)
  --catalog			Use the /v1/catalog endpoint to register
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *ServiceRegisterCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
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
	flags.StringVar(&c.checkHttp, "check-http", "", "")
	flags.StringVar(&c.checkScript, "check-script", "", "")
	flags.StringVar(&c.checkTtl, "check-ttl", "", "")
	flags.StringVar(&c.checkInterval, "check-interval", "", "")
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

	checkCount := 0
	if c.checkHttp != "" { checkCount = checkCount + 1}
	if c.checkScript != "" { checkCount = checkCount + 1}
	if c.checkTtl != "" { checkCount = checkCount + 1}

	if checkCount > 1 {
		c.UI.Error("Only one of --check-http, --check-script or --check-ttl can")
		c.UI.Error("be specified")
		return 1
	}

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
	}

	if checkCount == 1 {
		service.Check = &consulapi.AgentServiceCheck{
			Script:		c.checkScript,
			HTTP:		c.checkHttp,
			TTL:		c.checkTtl,
			Interval:	c.checkInterval,
		}
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
