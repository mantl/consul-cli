package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Watch functions

func newKvWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch <path>",
		Short: "Watch for changes to a K/V path",
		Long:  "Watch for changes to a K/V path",
		RunE:  kvWatch,
	}

	cmd.Flags().String("fields", "all", "Comma separated list of fields to return.")
	cmd.Flags().String("format", "prettyjson", "Output format. Supported options: text, json, prettyjson")
	cmd.Flags().String("delimited", "", "Output field delimiter")
	cmd.Flags().Bool("header", false, "Output a header row for text format")

	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)
	addWaitIndexOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func kvWatch(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

RETRY:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		queryOpts.WaitIndex = meta.LastIndex
		goto RETRY
	}

	if viper.GetString("template") != "" {
		return output(kv)
	} else {
		return outputKv(&kv)
	}
}

