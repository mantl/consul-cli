package commands

import (
	"github.com/spf13/cobra"
)

type CheckWarnOptions struct {
	Note string
}

func (c *Check) AddWarnSub(cmd *cobra.Command) {
	cfo := &CheckWarnOptions{}

	warnCmd := &cobra.Command{
		Use:   "warn <checkId>",
		Short: "Mark a local check as warning",
		Long:  "Mark a local check as warning",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Warn(args, cfo)
		},
	}

	oldWarnCmd := &cobra.Command{
		Use:        "check-warn <checkId>",
		Short:      "Mark a local check as warning",
		Long:       "Mark a local check as warning",
		Deprecated: "Use check warn",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Warn(args, cfo)
		},
	}

	warnCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")
	oldWarnCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")

	cmd.AddCommand(warnCmd)

	c.AddCommand(oldWarnCmd)
}

func (c *Check) Warn(args []string, cfo *CheckWarnOptions) error {
	if err := c.CheckIdArg(args); err != nil {
		return err
	}
	checkId := args[0]

	client, err := c.Agent()
	if err != nil {
		return err
	}

	err = client.WarnTTL(checkId, cfo.Note)
	if err != nil {
		return err
	}

	return nil
}
