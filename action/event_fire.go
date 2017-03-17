package action

import (
	"flag"
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

type eventFire struct {
	node string
	payload string
	service string
	tag string

	*config
}

func EventFireAction() Action {
	return &eventFire{
		config: &gConfig,
	}
}

func (e *eventFire) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	e.addDatacenterFlag(f)
	e.addOutputFlags(f, false)
	e.addRawFlag(f)

	f.StringVar(&e.node, "node", "", "Filter by node name")
	f.StringVar(&e.payload, "payload", "", "Event payload")
	f.StringVar(&e.service, "service", "", "Filter by service")
	f.StringVar(&e.tag, "tag", "", "Filter by service tag")

	return f
}

func (e *eventFire) Run(args []string) error {
	var event consulapi.UserEvent

	if e.raw.isSet() {
		if err := e.raw.readJSON(&event); err != nil {
			return err
		}
	} else {
		if len(args) != 1 {
			return fmt.Errorf("An event name must be specified")
		}
		eventName := args[0]

		var payload []byte
		if e.payload != "" {
			payload = []byte(e.payload)
		}

		event = consulapi.UserEvent{
			Name:          eventName,
			NodeFilter:    e.node,
			ServiceFilter: e.service,
			TagFilter:     e.tag,
			Payload:       payload,
		}
	}

	client, err := e.newEvent()
	if err != nil {
		return err
	}

	writeOpts := e.writeOptions()
	rval, _, err := client.Fire(&event, writeOpts)
	if err != nil {
		return err
	}

	return e.Output(rval)
}
