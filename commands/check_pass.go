package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Pass functions

func newCheckPassCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pass <checkId>",
		Short: "Mark a local check as passing",
		Long:  "Mark a local check as passing",
		RunE:  checkPass,
	}

	cmd.Flags().String("note", "", "Message to associate with check status")

	return cmd
}

func checkPass(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.PassTTL(checkId, viper.GetString("note"))
	if err != nil {
		return err
	}

	return nil
}
