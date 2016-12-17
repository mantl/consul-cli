package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Read functions

func newKvReadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <path>",
		Short: "Read a value from a given path",
		Long:  "Read a value from a given path",
		RunE:  kvRead,
	}

	cmd.Flags().String("fields", "value", "Comma separated list of fields to return")
	cmd.Flags().String("format", "text", "Output format. Supported options: text, json, prettyjson")
	cmd.Flags().String("delimiter", " ", "Output field delimiter")
	cmd.Flags().Bool("header", false, "Output a header row for text format")
	cmd.Flags().Bool("recurse", false, "Perform a recursive read")

	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func kvRead(cmd *cobra.Command, args []string) error {
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

	if viper.GetBool("recurse") {
		kvlist, _, err := client.List(path, queryOpts)
		if err != nil {
			return err
		}

		if kvlist == nil {
			return nil
		}

		if viper.GetString("template") != "" {
			return output(kvlist)
		} 

		return outputKv(&kvlist)
	} else {
		kv, _, err := client.Get(path, queryOpts)
		if err != nil {
			return err
		}

		if kv == nil {
			return nil
		}

		if viper.GetString("template") != "" {
			return output(kv)
		}

		return outputKv(kv)
	}
}

