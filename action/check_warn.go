package action

import (
	"flag"
	"fmt"
)

type checkWarn struct {
	note string

	*config
}

func CheckWarnAction() Action {
	return &checkWarn{
		config: &gConfig,
	}
}

func (c *checkWarn) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&c.note, "note", "", "Message to associate with check status")

	return f
}

func (c *checkWarn) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.WarnTTL(checkId, c.note)
	if err != nil {
		return err
	}

	return nil
}
