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
	f := newFlagSet()

	c.addDatacenterFlag(f)
	c.addConsistencyFlags(f)
	c.addOutputFlags(f, false)

	return f
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
