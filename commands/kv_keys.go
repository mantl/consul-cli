package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Keys functions

func newKvKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys <path>",
		Short: "List K/V keys",
		Long:  "List K/V keys",
		RunE:  kvKeys,
	}

	cmd.Flags().String("separator", "", "List keys only up to a given separator")
	addDatacenterOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func kvKeys(cmd *cobra.Command, args []string) error {
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
	data, _, err := client.Keys(path, viper.GetString("separator"), queryOpts)
	if err != nil {
		return err
	}

	viper.Set("template", kv_outputTemplate)

	return output(data)
}

var kv_outputTemplate = `{{range .}}{{.}}
{{end}}`

