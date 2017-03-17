package action

import (
	"flag"
)

type coordNodes struct {
	*config
}

func CoordNodesAction() Action {
	return &coordNodes{
		config: &gConfig,
	}
}

func (c *coordNodes) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	c.addDatacenterFlag(f)
	c.addOutputFlags(f, false)
	c.addConsistencyFlags(f)

	return f
}

func (c *coordNodes) Run(args []string) error {
	client, err := c.newCoordinate()
	if err != nil {
		return err
	}

	queryOpts := c.queryOptions()
	data, _, err := client.Nodes(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(data)
}
