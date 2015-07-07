package command

import (
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type KVDeleteCommand struct {
	Meta
}

func (c *KVDeleteCommand) Help() string {
	helpText := `
Usage: consul-cli kv-delete [options] path

  Delete a given path from the K/V

Options:

` + c.ConsulHelp() +
`  --modifyindex=<ModifyIndex>	Perform a Check-and-Set delete
				(default: not set)
  --recurse			Perform a recursive delete
				(default: false)
`

	return strings.TrimSpace(helpText)
}

func (c *KVDeleteCommand) Run (args[]string) int {
	var modifyIndex string
	var doRecurse bool

	flags := c.Meta.FlagSet()
	flags.StringVar(&modifyIndex, "modifyindex", "", "")
	flags.BoolVar(&doRecurse, "recurse", false, "")
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

	writeOpts := c.WriteOptions()

	switch {
	case doRecurse:
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	case modifyIndex != "":
		m, err := strconv.ParseUint(modifyIndex, 0, 64)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
		kv := consulapi.KVPair{
			Key:		path,
			ModifyIndex:	m,
		}

		success, _, err := client.DeleteCAS(&kv, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if !success {
			return 1
		}
	default:
		_, err := client.Delete(path, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	}

	return 0
}

func (c *KVDeleteCommand) Synopsis() string {
	return "Delete a path"
}
