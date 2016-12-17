package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// State functions

func newHealthStateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Get the checks in a given state",
		Long:  "Get the checks in a given state",
		RunE:  healthState,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthState(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Check state must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one check state allowed")
	}
	state := strings.ToLower(args[0])

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	s, _, err := client.State(state, queryOpts)
	if err != nil {
		return err
	}

	return output(s)
}
