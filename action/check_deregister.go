package action

import (
	"flag"
	"fmt"
)

type checkDeregister struct {
	*config
}

func CheckDeregisterAction() Action {
	return &checkDeregister{
		config: &gConfig,
	}
}

func (c *checkDeregister) CommandFlags() *flag.FlagSet {
	return c.newFlagSet(FLAG_NONE)
}

func (c *checkDeregister) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.CheckDeregister(checkId)
	if err != nil {
		return err
	}

	return nil
}
