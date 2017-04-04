package action

import (
	"flag"
	"fmt"
)

type catalogService struct {
	tag string

	*config
}

func CatalogServiceAction() Action {
	return &catalogService{
		config: &gConfig,
	}
}

func (c *catalogService) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&c.tag, "tag", "", "Service tag to filter on")

	c.addDatacenterFlag(f)
	c.addOutputFlags(f, false)
	c.addConsistencyFlags(f)

	return f
}

func (c *catalogService) Run(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service allowed")
	}

	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	queryOpts := c.queryOptions()
	config, _, err := client.Service(args[0], c.tag, queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
