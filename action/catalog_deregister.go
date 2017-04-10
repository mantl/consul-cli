package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type catalogDeregister struct {
	node      string
	serviceId string
	checkId   string

	*config
}

func CatalogDeregisterAction() Action {
	return &catalogDeregister{
		config: &gConfig,
	}
}

func (c *catalogDeregister) CommandFlags() *flag.FlagSet {
	f := c.newFlagSet(FLAG_DATACENTER)

	f.StringVar(&c.node, "node", "", "Consul node name. Required")
	f.StringVar(&c.serviceId, "service-id", "", "Service ID to deregister")
	f.StringVar(&c.checkId, "check-id", "", "Check ID to deregister")

	return f
}

func (c *catalogDeregister) Run(args []string) error {
	client, err := c.newCatalog()
	if err != nil {
		return err
	}

	if c.node == "" {
		return fmt.Errorf("Node name is required for catalog deregistration")
	}

	writeOpts := c.writeOptions()
	_, err = client.Deregister(&consulapi.CatalogDeregistration{
		Node:       c.node,
		Datacenter: c.dc,
		ServiceID:  c.serviceId,
		CheckID:    c.checkId,
	}, writeOpts)

	return err
}
