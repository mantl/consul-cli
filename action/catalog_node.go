package action

import (
	"flag"
	"fmt"
)

type catalogNode struct {
	*config
}

func CatalogNodeAction() Action {
	return &catalogNode{
		config: &gConfig,
	}
}

func (c *catalogNode) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	c.addDatacenterFlag(f)
	c.addConsistencyFlags(f)
	c.addOutputFlags(f, false)

	return f
}

func (c *catalogNode) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	queryOpts := c.queryOptions()
	config, _, err := client.Node(args[0], queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
