package command

import (
	"strings"
)

type CatalogNodesCommand struct {
	Meta
}

func (c *CatalogNodesCommand) Help() string {
	helpText := `
Usage: consul-cli catalog-nodes [options]

  Get all the nodes registered with a given DC

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *CatalogNodesCommand) Run(args []string) int {
	c.AddDataCenter()
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
	queryOpts := c.QueryOptions()
        config, _, err := client.Nodes(queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *CatalogNodesCommand) Synopsis() string {
	return "Get all the nodes registered with a given DC"
}
