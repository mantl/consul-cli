package action

import (
	"flag"
	"fmt"
)

type checkPass struct {
	note string

	*config
}

func CheckPassAction() Action {
	return &checkPass{
		config: &gConfig,
	}
}

func (c *checkPass) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&c.note, "note", "", "Message to associate with check status")

	return f
}

func (c *checkPass) Run(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	client, err := c.newAgent()
	if err != nil {
		return err
	}

	err = client.PassTTL(checkId, c.note)
	if err != nil {
		return err
	}

	return nil
}
