package action

import (
	"flag"
)

type catalogServices struct {
	*config
}

func CatalogServicesAction() Action {
	return &catalogServices{
		config: &gConfig,
	}
}

func (c *catalogServices) CommandFlags() *flag.FlagSet {
	return newFlagSet()
}

//	addDatacenterOption(cmd)
//	addTemplateOption(cmd)
//	addConsistencyOptions(cmd)

func (c *catalogServices) Run(args []string) error {
	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	queryOpts := c.queryOptions()
	config, _, err := client.Services(queryOpts)
	if err != nil {
		return err
	}

	return c.Output(config)
}
