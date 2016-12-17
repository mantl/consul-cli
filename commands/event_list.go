package commands

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// List functions

func newEventListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Lists the most recent events the agent has seen",
		Long:  "Lists the most recent events the agent has seen",
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
