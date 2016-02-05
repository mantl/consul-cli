package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

func (a *Agent) AddForceLeaveSub(c *cobra.Command) {
	forceLeaveCmd := &cobra.Command{
		Use:   "force-leave <node name>",
		Short: "Force the removal of a node",
		Long:  "Force the removal of a node",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.ForceLeave(args)
		},
	}

	oldForceLeaveCmd := &cobra.Command{
		Use:        "agent-force-leave",
		Short:      "Force the removal of a node",
		Long:       "Force the removal of a node",
		Deprecated: "Use agent force-leave",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.ForceLeave(args)
		},
	}

	a.AddTemplateOption(forceLeaveCmd)
	c.AddCommand(forceLeaveCmd)

	a.AddCommand(oldForceLeaveCmd)
}

func (a *Agent) ForceLeave(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Name not provided")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	client, err := a.Agent()
	if err != nil {
		return err
	}

	return client.ForceLeave(args[0])
}
