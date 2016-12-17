package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Fire functions

func newEventFireCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fire <name>",
		Short: "Fires a new user event",
		Long:  "Fires a new user event",
		RunE:  eventFire,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addRawOption(cmd)

	cmd.Flags().String("node", "", "Filter by node name")
	cmd.Flags().String("payload", "", "Event payload")
	cmd.Flags().String("service", "", "Filter by service")
	cmd.Flags().String("tag", "", "Filter by service tag")

	return cmd
}

func eventFire(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var event consulapi.UserEvent

	if raw := viper.GetString("raw"); raw != "" {
		if err := readRawJSON(raw, &event); err != nil {
			return err
		}
	} else {
		if len(args) != 1 {
			return fmt.Errorf("An event name must be specified")
		}
		eventName := args[0]

		var payload []byte

		if ps := viper.GetString("payload"); ps != "" {
			payload = []byte(ps)
		}

		event = consulapi.UserEvent{
			Name:          eventName,
			NodeFilter:    viper.GetString("node"),
			ServiceFilter: viper.GetString("service"),
			TagFilter:     viper.GetString("tag"),
			Payload:       payload,
		}
	}

	client, err := newEvent()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()
	rval, _, err := client.Fire(&event, writeOpts)
	if err != nil {
		return err
	}

	return output(rval)
}
