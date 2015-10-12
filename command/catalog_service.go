package command

import (
	"strings"
)

type CatalogServiceCommand struct {
	Meta
        tag             string
}

func (c *CatalogServiceCommand) Help() string {
	helpText := `
Usage: consul-cli catalog-services [options]

  Get the nodes providing a service

Options:
` + c.ConsulHelp() +
`  --tag                         Service tag to filter on
                                (default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *CatalogServiceCommand) Run(args []string) int {
        c.AddDataCenter()
        flags := c.Meta.FlagSet()
        flags.StringVar(&c.tag, "tag", "", "")
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
        service := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Catalog()
	queryOpts := c.QueryOptions()
        config, _, err := client.Service(service, c.tag, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *CatalogServiceCommand) Synopsis() string {
	return "Get the nodes providing a service"
}
