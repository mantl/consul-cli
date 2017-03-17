package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type catalogRegister struct {
	node string

	*config
}

func CatalogRegisterAction() Action {
	return &catalogRegister{
		config: &gConfig,
	}
}

func (c *catalogRegister) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&c.node, "node", "", "Service node")

	c.addDatacenterFlag(f)
	c.addServiceFlags(f)

	return f
}

func (c *catalogRegister) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single service name must be specified")
	}
	service := args[0]

	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	writeOpts := c.writeOptions()

	_, err = client.Register(&consulapi.CatalogRegistration{
		Node:       c.node,
		Address:    c.service.address,
		Datacenter: c.dc,
		Service: &consulapi.AgentService{
			ID:                c.service.id,
			Service:           service,
			Tags:              c.service.tags,
			Port:              c.service.port,
			EnableTagOverride: c.service.overrideTag,
		},
	},
		writeOpts)

	return err
}
