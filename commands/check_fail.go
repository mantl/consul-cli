package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Fail functions

func newCheckFailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fail <checkId>",
		Short: "Mark a local check as critical",
		Long:  "Mark a local check as critical",
		RunE:  checkFail,
	}

	cmd.Flags().String("note", "", "Message to associate with check status")

	return cmd
}

func checkFail(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.FailTTL(checkId, viper.GetString("note"))
	if err != nil {
		return err
	}

	return nil
}
