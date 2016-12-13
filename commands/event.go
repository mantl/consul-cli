package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newEventCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "event",
		Short: "Consul /event endpoint interface",
		Long:  "Consul /event endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newEventFireCommand())
	cmd.AddCommand(newEventListCommand())

	return cmd
}

// fire functions

func newEventFireCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "fire <name>",
		Short: "Fires a new user event",
		Long: "Fires a new user event",
		RunE: eventFire,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	cmd.Flags().String("node", "", "Filter by node name")
	cmd.Flags().String("payload", "", "Event payload")
	cmd.Flags().String("service", "", "Filter by service")
	cmd.Flags().String("tag", "", "Filter by service tag")

	return cmd
}

func eventFire(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	if len(args) != 1 {
		return fmt.Errorf("An event name must be specified")
	}

	client, err := newEvent()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	var payload []byte

	if ps := viper.GetString("payload"); ps != "" {
		payload = []byte(ps)
	}

	rval, _, err := client.Fire(
		&consulapi.UserEvent{
			Name: args[0],
			NodeFilter: viper.GetString("node"),
			ServiceFilter: viper.GetString("service"),
			TagFilter: viper.GetString("tag"),
			Payload: payload,
		},
		writeOpts,
		)
	if err != nil {
		return err
	}

	return output(rval)
}

// list functions

func newEventListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the most recent events the agent has seen",
		Long: "Lists the most recent events the agent has seen",
		RunE:  eventList,
	}

	cmd.Flags().String("name", "", "Event name to filter on")

	addTemplateOption(cmd)

	return cmd
}

func eventList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newEvent()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	data, _, err := client.List(viper.GetString("name"), queryOpts)
	if err != nil {
		return err
	}

	return output(data)
}
