package command

import (
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type CheckRegisterCommand struct {
	Meta
	id			string
	script			string
	http			string
	ttl			string
	interval		string
	notes			string
	serviceID		string
}

func (c *CheckRegisterCommand) Help() string {
	helpText := `
Usage: consul-cli check-register [options] checkName

  Register a new local check

  If --id is not specified, the checkName is used. There cannot
be duplicate IDs per agent however.

  Only one of --http, --script and --ttl can be specified.

Options:
` + c.ConsulHelp() + 
`  --id				Service Id
				(default: not set)
  --http			A URL to GET every interval
				(default: not set)
  --script			A script to run every interval
				(default: not set)
  --ttl				Fail if TTL expires before service 
				checks in
				(default: not set)
  --interval			Interval between checks
				(default: not set)
  --service-id			Service ID to associate check
				(default: not set)
  --notes			Description of the check
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *CheckRegisterCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.StringVar(&c.id, "id", "", "")
	flags.StringVar(&c.http, "http", "", "")
	flags.StringVar(&c.script, "script", "", "")
	flags.StringVar(&c.ttl, "ttl", "", "")
	flags.StringVar(&c.interval, "interval", "", "")
	flags.StringVar(&c.notes, "notes", "", "")
	flags.StringVar(&c.serviceID, "service-id", "", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Check name must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}
	checkName := extra[0]

	checkCount := 0
	if c.http != "" { checkCount = checkCount + 1}
	if c.script != "" { checkCount = checkCount + 1}
	if c.ttl != "" { checkCount = checkCount + 1}

	if checkCount > 1 {
		c.UI.Error("Only one of --http, --script or --ttl can be specified")
		return 1
	}

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	check := &consulapi.AgentCheckRegistration{
		ID:			c.id,
		Name:			checkName,
		ServiceID:		c.serviceID,
		Notes:			c.notes,
		AgentServiceCheck:	consulapi.AgentServiceCheck{
			Script:		c.script,
			HTTP:		c.http,
			Interval:	c.interval,
			TTL:		c.ttl,
		},
	}

	client := consul.Agent()
	err = client.CheckRegister(check)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	return 0
}

func (c *CheckRegisterCommand) Synopsis() string {
	return "Register a new local check"
}
