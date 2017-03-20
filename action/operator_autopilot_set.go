// +build consul8
//

package action

import (
	"flag"
	"fmt"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

type operatorAutopilotSet struct {
	modifyIndex string
	cleanDead bool

	*config
}

func OperatorAutopilotSetAction() Action {
	return &operatorAutopilotSet{
		config: &gConfig,
	}
}

func (o *operatorAutopilotSet) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&o.modifyIndex, "modifyindex", "", "Perform a check-and-set operation")
	f.BoolVar(&o.cleanDead, "clean-dead-servers", false, "Remove dead servers automatically when a new server is added")

	o.addDatacenterFlag(f)

	return f
}

func (o *operatorAutopilotSet) Run(args []string) error {
	client, err := o.newOperator()
	if err != nil {
		return err
	}

	writeOpts := o.writeOptions()

	ac := &consulapi.AutopilotConfiguration{
		CleanupDeadServers: o.cleanDead,
	}

	if o.modifyIndex == "" {
		_, err := client.AutopilotSetConfiguration(ac, writeOpts)
		if err != nil {
			return err
		}
	} else {
		i, err := strconv.ParseUint(o.modifyIndex, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing modifyIndex: %v", o.modifyIndex)
		}
		ac.ModifyIndex = i

		success, _, err := client.AutopilotCASConfiguration(ac, writeOpts)
		if err != nil {
			return err
		}

		if !success {
			return fmt.Errorf("Failed to write Autopilot configuration")
		}
	}

	return nil
}
