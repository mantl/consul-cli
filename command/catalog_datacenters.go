package command

import (
	"strings"
)

type CatalogDatacentersCommand struct {
	Meta
}

func (c *CatalogDatacentersCommand) Help() string {
	helpText := `
Usage: consul-cli catalog-services [options]

  Get all the datacenters known by the Consul server


Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *CatalogDatacentersCommand) Run(args []string) int {
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Catalog()
        config, err := client.Datacenters()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *CatalogDatacentersCommand) Synopsis() string {
	return "Get all the datacenters known by the Consul server"
}
