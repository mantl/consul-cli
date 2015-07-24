package command

import (
	"strings"
)

type KVWatchCommand struct {
	Meta
	format		OutputFormat
	fieldsRaw	string
}

func (c *KVWatchCommand) Help() string {
	helpText := `
Usage: consul-cli kv-watch [options] path

  Watch for changes to a K/V path

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
`

	return strings.TrimSpace(helpText)
}

func (c *KVWatchCommand) Run(args []string) int {
	c.AddDataCenter()
	c.AddWait()

	flags := c.Meta.FlagSet()
	flags.StringVar(&c.fieldsRaw, "fields", "value", "")
	flags.StringVar(&c.format.Type, "format", "text", "")
	flags.StringVar(&c.format.Delimiter, "delimiter", " ", "")
	flags.BoolVar(&c.format.Header, "header", false, "")
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

	kvo := NewKVOutput(c.UI, c.fieldsRaw)

RETRY:
	c.UI.Output("Getting path")
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	if kv == nil {
		c.UI.Output("kv == nil, setting WaitIndex")
		queryOpts.WaitIndex = meta.LastIndex
		goto RETRY
	}

	kvo.Output(kv, c.format)

	return 0
}

func (c *KVWatchCommand) Synopsis() string {
	return "Watch a path for changes"
}
