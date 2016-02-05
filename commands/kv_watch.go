package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KvWatchOptions struct {
	Format    OutputFormat
	FieldsRaw string
}

func (k *Kv) AddWatchSub(cmd *cobra.Command) {
	kwo := &KvWatchOptions{}

	watchCmd := &cobra.Command{
		Use:   "watch <path>",
		Short: "Watch for changes to a K/V path",
		Long:  "Watch for changes to a K/V path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Watch(args, kwo)
		},
	}

	oldWatchCmd := &cobra.Command{
		Use:        "kv-watch <path>",
		Short:      "Watch for changes to a K/V path",
		Long:       "Watch for changes to a K/V path",
		Deprecated: "Use kv watch",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Watch(args, kwo)
		},
	}

	watchCmd.Flags().StringVar(&kwo.FieldsRaw, "fields", "", "Comma separated list of fields to return.")
	watchCmd.Flags().StringVar(&kwo.Format.Type, "format", "", "Output format. Supported options: text, json, prettyjson")
	watchCmd.Flags().StringVar(&kwo.Format.Delimiter, "delimited", "", "Output field delimiter")
	watchCmd.Flags().BoolVar(&kwo.Format.Header, "header", false, "Output a header row for text format")
	k.AddDatacenterOption(watchCmd)
	k.AddWaitIndexOption(watchCmd)
	k.AddTemplateOption(watchCmd)

	oldWatchCmd.Flags().StringVar(&kwo.FieldsRaw, "fields", "", "Comma separated list of fields to return.")
	oldWatchCmd.Flags().StringVar(&kwo.Format.Type, "format", "", "Output format. Supported options: text, json, prettyjson")
	oldWatchCmd.Flags().StringVar(&kwo.Format.Delimiter, "delimited", "", "Output field delimiter")
	oldWatchCmd.Flags().BoolVar(&kwo.Format.Header, "header", false, "Output a header row for text format")
	k.AddDatacenterOption(oldWatchCmd)
	k.AddWaitIndexOption(oldWatchCmd)

	cmd.AddCommand(watchCmd)

	k.AddCommand(oldWatchCmd)
}

func (k *Kv) Watch(args []string, kwo *KvWatchOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Key path must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one path allowed")
	}
	path := args[0]

	client, err := k.KV()
	if err != nil {
		return err
	}

	queryOpts := k.QueryOptions()

	kvo := NewKVOutput(k.Out, k.Err, kwo.FieldsRaw)

RETRY:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		queryOpts.WaitIndex = meta.LastIndex
		goto RETRY
	}

	if k.Template != "" {
		return k.Output(kv)
	} else {
		return kvo.Output(kv, kwo.Format)
	}
}
