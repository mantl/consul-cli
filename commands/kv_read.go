package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

type KvReadOptions struct {
	Format    OutputFormat
	FieldsRaw string
	Recurse   bool
}

func (k *Kv) AddReadSub(cmd *cobra.Command) {
	kro := &KvReadOptions{}

	readCmd := &cobra.Command{
		Use:   "read <path>",
		Short: "Read a value from a given path",
		Long:  "Read a value from a given path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Read(args, kro)
		},
	}

	oldReadCmd := &cobra.Command{
		Use:        "kv-read <path>",
		Short:      "Read a value from a given path",
		Long:       "Read a value from a given path",
		Deprecated: "Use kv read",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Read(args, kro)
		},
	}

	readCmd.Flags().StringVar(&kro.FieldsRaw, "fields", "value", "Comma separated list of fields to return")
	readCmd.Flags().StringVar(&kro.Format.Type, "format", "text", "Output format. Supported options: text, json, prettyjson")
	readCmd.Flags().StringVar(&kro.Format.Delimiter, "delimiter", " ", "Output field delimiter")
	readCmd.Flags().BoolVar(&kro.Format.Header, "header", false, "Output a header row for text format")
	readCmd.Flags().BoolVar(&kro.Recurse, "recurse", false, "Perform a recursive read")
	k.AddDatacenterOption(readCmd)
	k.AddTemplateOption(readCmd)

	oldReadCmd.Flags().StringVar(&kro.FieldsRaw, "fields", "value", "Comma separated list of fields to return")
	oldReadCmd.Flags().StringVar(&kro.Format.Type, "format", "text", "Output format. Supported options: text, json, prettyjson")
	oldReadCmd.Flags().StringVar(&kro.Format.Delimiter, "delimiter", " ", "Output field delimiter")
	oldReadCmd.Flags().BoolVar(&kro.Format.Header, "header", false, "Output a header row for text format")
	oldReadCmd.Flags().BoolVar(&kro.Recurse, "recurse", false, "Perform a recursive oldRead")
	k.AddDatacenterOption(oldReadCmd)

	cmd.AddCommand(readCmd)

	k.AddCommand(oldReadCmd)
}

func (k *Kv) Read(args []string, kro *KvReadOptions) error {
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

	kvo := NewKVOutput(k.Out, k.Err, kro.FieldsRaw)

	if kro.Recurse {
		kvlist, _, err := client.List(path, queryOpts)
		if err != nil {
			return err
		}

		if kvlist == nil {
			return nil
		}

		if k.Template != "" {
			return k.Output(kvlist)
		} else {
			return kvo.OutputList(&kvlist, kro.Format)
		}
	} else {
		kv, _, err := client.Get(path, queryOpts)
		if err != nil {
			return err
		}

		if kv == nil {
			return nil
		}

		if k.Template != "" {
			return k.Output(kv)
		} else {
			return kvo.Output(kv, kro.Format)
		}
	}

	return nil
}
