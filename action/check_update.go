package action

import (
	"flag"
	"fmt"
)

type checkUpdate struct {
	status string
	output string

	*config
}

func CheckUpdateAction() Action {
	return &checkUpdate{
		config: &gConfig,
	}
}

func (c *checkUpdate) CommandFlags() *flag.FlagSet {
	f := c.newFlagSet(FLAG_NONE)

	f.StringVar(&c.status, "status", "", "Check status. One of passing, warning or critical")
	f.StringVar(&c.output, "output", "", "Optional human readable message with the status of the check")

	return f
}

func (c *checkUpdate) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.UpdateTTL(checkId, c.output, c.status)
	if err != nil {
		return err
	}

	return nil
}
