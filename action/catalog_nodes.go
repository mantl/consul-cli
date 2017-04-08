package action

import (
	"flag"
)

type catalogNodes struct {
	*config
}

func CatalogNodesAction() Action {
	return &catalogNodes{
		config: &gConfig,
	}
}

func (c *catalogNodes) CommandFlags() *flag.FlagSet {
	return c.newFlagSet(FLAG_DATACENTER, FLAG_CONSISTENCY, FLAG_OUTPUT)
}

func (c *catalogNodes) Run(args []string) error {
	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	queryOpts := c.queryOptions()
	config, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
