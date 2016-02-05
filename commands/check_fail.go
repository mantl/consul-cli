package commands

import (
	"github.com/spf13/cobra"
)

type CheckFailOptions struct {
	Note string
}

func (c *Check) AddFailSub(cmd *cobra.Command) {
	cfo := &CheckFailOptions{}

	failCmd := &cobra.Command{
		Use:   "fail <checkId>",
		Short: "Mark a local check as critical",
		Long:  "Mark a local check as critical",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Fail(args, cfo)
		},
	}

	oldFailCmd := &cobra.Command{
		Use:        "check-fail <checkId>",
		Short:      "Mark a local check as critical",
		Long:       "Mark a local check as critical",
		Deprecated: "Use check fail",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Fail(args, cfo)
		},
	}

	failCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")
	oldFailCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")

	cmd.AddCommand(failCmd)

	c.AddCommand(oldFailCmd)
}

func (c *Check) Fail(args []string, cfo *CheckFailOptions) error {
	if err := c.CheckIdArg(args); err != nil {
		return err
	}
	checkId := args[0]

	client, err := c.Agent()
	if err != nil {
		return err
	}

	err = client.FailTTL(checkId, cfo.Note)
	if err != nil {
		return err
	}

	return nil
}
