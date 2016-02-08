package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KvKeysOptions struct {
	Separator string
}

func (k *Kv) AddKeysSub(cmd *cobra.Command) {
	kko := &KvKeysOptions{}

	keysCmd := &cobra.Command{
		Use:   "keys <path>",
		Short: "List K/V keys",
		Long: "List K/V keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Keys(args, kko)
		},
	}

	keysCmd.Flags().StringVar(&kko.Separator, "separator", "", "List keys only up to a given separator")
	k.AddDatacenterOption(keysCmd)

	cmd.AddCommand(keysCmd)
}

func (k *Kv) Keys(args []string, kko *KvKeysOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Key path must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one key path allowed")
	}
	path := args[0]

	client, err := k.KV()
	if err != nil {
		return err
	}

	queryOpts := k.QueryOptions()
	data, _, err := client.Keys(path, kko.Separator, queryOpts)
	if err != nil {
		return err
	}

	k.Template = outputTemplate
	k.Output(data)

	return nil
}

var outputTemplate = `{{range .}}{{.}}
{{end}}`
