package commands

import (
	"github.com/spf13/cobra"
)

func (c *Check) AddDeregisterSub(cmd *cobra.Command) {
	deregisterCmd := &cobra.Command{
		Use: "deregister",
		Short: "Remove a check from the agent",
		Long: "Remove a check from the agent",
		RunE: func(cmd *cobra.Command, args[]string) error {
			return c.Deregister(args)
		},
	}

	oldDeregisterCmd := &cobra.Command{
		Use: "deregister",
		Short: "Remove a check from the agent",
		Long: "Remove a check from the agent",
		Deprecated: "Use acl deregister",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args[]string) error {
			return c.Deregister(args)
		},
	}

	cmd.AddCommand(deregisterCmd)

	c.AddCommand(oldDeregisterCmd)
}

func (c *Check) Deregister(args []string) error {
	if err := c.CheckIdArg(args); err != nil {
		return err
	}
	checkId := args[0]

	consul, err := c.Client()
	if err != nil {	
		return err
	}

	client := consul.Agent()
	err = client.CheckDeregister(checkId)
	if err != nil {
		return err
	}

	return nil
}
