package command

import (
	"strings"
)

type KVReadCommand struct {
	Meta
}

func (c *KVReadCommand) Help() string {
	helpText := `
Usage: consul-cli kv-read [options] path

  Read a value from a given path.

Options:

` + c.ConsulHelp() +
`  --fields=<f1,f2,...>		Comma separated list of fields to return.
				(default: value)
  --format=text			Output format. Supported options: text, json, prettyjson
				(default: text)
  --delimiter=			Output field delimiter.
				(default: " ")
  --header			Output a header row for text format
				(default: false)
  --recurse			Perform a recursive read
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *KVReadCommand) Run(args []string) int {
	var format OutputFormat
	var fieldsRaw string
	var recurse bool

	flags := c.Meta.FlagSet(true)
	flags.StringVar(&fieldsRaw, "fields", "value", "")
	flags.StringVar(&format.Type, "format", "text", "")
	flags.StringVar(&format.Delimiter, "delimiter", " ", "")
	flags.BoolVar(&format.Header, "header", false, "")
	flags.BoolVar(&recurse, "recurse", false, "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 1 {
		c.UI.Error("Key path must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	path := extra[0]

	consul, err := c.Client()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	queryOpts := c.QueryOptions()

	kvo := NewKVOutput(c.UI, fieldsRaw)

	if recurse {
		kvlist, _, err := client.List(path, queryOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if kvlist == nil {
			return 0
		}

		kvo.OutputList(&kvlist, format)
	} else {
		kv, _, err := client.Get(path, queryOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if kv == nil {
			return 0
		}

		kvo.Output(kv, format)
	}

	return 0
}

func (c *KVReadCommand) Synopsis() string {
	return "Read a value"
}
