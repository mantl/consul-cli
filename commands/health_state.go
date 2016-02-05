package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func (h *Health) AddStateSub(cmd *cobra.Command) {
	stateCmd := &cobra.Command{
		Use:   "state",
		Short: "Get the checks in a given state",
		Long:  "Get the checks in a given state",
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.State(args)
		},
	}

	oldStateCmd := &cobra.Command{
		Use:        "state",
		Short:      "Get the checks in a given state",
		Long:       "Get the checks in a given state",
		Deprecated: "Use health state",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return h.State(args)
		},
	}

	h.AddDatacenterOption(stateCmd)
	h.AddTemplateOption(stateCmd)
	h.AddDatacenterOption(oldStateCmd)

	cmd.AddCommand(stateCmd)

	h.AddCommand(oldStateCmd)
}

func (h *Health) State(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Check state must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one check state allowed")
	}
	state := strings.ToLower(args[0])

	client, err := h.Health()
	if err != nil {
		return err
	}

	queryOpts := h.QueryOptions()

	s, _, err := client.State(state, queryOpts)
	if err != nil {
		return err
	}

	return h.Output(s)
}
