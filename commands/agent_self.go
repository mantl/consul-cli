package commands

import (
	"github.com/spf13/cobra"
)

func (a *Agent) AddSelfSub(c *cobra.Command) {
	selfCmd := &cobra.Command{
		Use: "self",
		Short: "Get agent configuration",
		Long: "Get agent configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Self(args)
		},
	}

	oldSelfCmd := &cobra.Command{
		Use: "agent-self",
		Short: "Get agent configuration",
		Long: "Get agent configuration",
		Deprecated: "Use agent self",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return a.Self(args)
		},
	}

	a.AddTemplateOption(selfCmd)
	c.AddCommand(selfCmd)

	a.AddCommand(oldSelfCmd)
}

func (a *Agent) Self(args []string) error {
	consul, err := a.Client()
	if err != nil {
		return err
	}

	client := consul.Agent()
	config, err := client.Self()
	if err != nil {
		return err
	}

	return a.Output(config)
}
