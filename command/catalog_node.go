package command

import (
	"strings"
)

type CatalogNodeCommand struct {
	Meta
}

func (c *CatalogNodeCommand) Help() string {
	helpText := `
Usage: consul-cli catalog-node [options]

  Get the services provided by a node

Options:
` + c.ConsulHelp()

	return strings.TrimSpace(helpText)
}

func (c *CatalogNodeCommand) Run(args []string) int {
        c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

        extra := flags.Args()
        if len(extra) < 1 {
                c.UI.Error("Node name must be specified")
                c.UI.Error("")
                c.UI.Error(c.Help())
                return 1
        }
        node := extra[0]

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}

	client := consul.Catalog()
	queryOpts := c.QueryOptions()
        config, _, err := client.Node(node, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.OutputJSON(config, true)

	return 0
}

func (c *CatalogNodeCommand) Synopsis() string {
	return "Get the services provided by a node"
}
