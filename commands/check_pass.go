package commands

import (
	"github.com/spf13/cobra"
)

type CheckPassOptions struct {
	Note string
}

func (c *Check) AddPassSub(cmd *cobra.Command) {
	cfo := &CheckPassOptions{}

	passCmd := &cobra.Command{
		Use:   "pass <checkId>",
		Short: "Mark a local check as passing",
		Long:  "Mark a local check as passing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Pass(args, cfo)
		},
	}

	oldPassCmd := &cobra.Command{
		Use:        "check-pass <checkId>",
		Short:      "Mark a local check as passing",
		Long:       "Mark a local check as passing",
		Deprecated: "Use check pass",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Pass(args, cfo)
		},
	}

	passCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")
	oldPassCmd.Flags().StringVar(&cfo.Note, "note", "", "Message to associate with check status")

	cmd.AddCommand(passCmd)

	c.AddCommand(oldPassCmd)
}

func (c *Check) Pass(args []string, cfo *CheckPassOptions) error {
	if err := c.CheckIdArg(args); err != nil {
		return err
	}
	checkId := args[0]

	client, err := c.Agent()
	if err != nil {
		return err
	}

	err = client.PassTTL(checkId, cfo.Note)
	if err != nil {
		return err
	}

	return nil
}
