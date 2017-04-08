package action

import (
	"flag"
)

type coordDatacenters struct {
	*config
}

func CoordDatacentersAction() Action {
	return &coordDatacenters{
		config: &gConfig,
	}
}

func (c *coordDatacenters) CommandFlags() *flag.FlagSet {
	return c.newFlagSet(FLAG_OUTPUT)
}

func (c *coordDatacenters) Run(args []string) error {
	client, err := c.newCoordinate()
	if err != nil {
		return err
	}

	data, err := client.Datacenters()
	if err != nil {
		return err
	}

	return c.Output(data)
}
