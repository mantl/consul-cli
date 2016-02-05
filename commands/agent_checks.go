package commands

import (
	"github.com/spf13/cobra"
)

func (a *Agent) AddChecksSub(c *cobra.Command) {
	checksCmd := &cobra.Command{
		Use:   "checks",
		Short: "Get the checks the agent is managing",
		Long:  "Get the checks the agent is managing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Checks(args)
		},
	}

	oldChecksCmd := &cobra.Command{
		Use:        "agent-checks",
		Short:      "Get the checks the agent is managing",
		Long:       "Get the checks the agent is managing",
		Deprecated: "Use agent checks",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Checks(args)
		},
	}

	a.AddTemplateOption(checksCmd)
	c.AddCommand(checksCmd)

	a.AddCommand(oldChecksCmd)
}

func (a *Agent) Checks(args []string) error {
	client, err := a.Agent()
	if err != nil {
		return err
	}

	config, err := client.Checks()
	if err != nil {
		return err
	}

	return a.Output(config)
}
