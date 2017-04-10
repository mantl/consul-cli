package action

import (
	"flag"
)

type eventList struct {
	name string

	*config
}

func EventListAction() Action {
	return &eventList{
		config: &gConfig,
	}
}

func (e *eventList) CommandFlags() *flag.FlagSet {
	f := e.newFlagSet(FLAG_OUTPUT, FLAG_BLOCKING)

	f.StringVar(&e.name, "name", "", "Event name to filter on")

	return f
}

func (e *eventList) Run(args []string) error {
	client, err := e.newEvent()
	if err != nil {
		return err
	}

	queryOpts := e.queryOptions()

	data, _, err := client.List(e.name, queryOpts)
	if err != nil {
		return err
	}

	return e.Output(data)
}
