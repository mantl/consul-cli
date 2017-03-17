package action

import (
	"flag"
	"fmt"
)

type checkFail struct {
	note string

	*config
}

func CheckFailAction() Action {
	return &checkFail{
		config: &gConfig,
	}
}

func (c *checkFail) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&c.note, "note", "", "Message to associate with check status")

	return f
}

func (c *checkFail) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.FailTTL(checkId, c.note)
	if err != nil {
		return err
	}

	return nil
}
