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
	return c.newFlagSet(FLAG_DATACENTER, FLAG_OUTPUT, FLAG_CONSISTENCY)
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
