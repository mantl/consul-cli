package command

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type KVWriteCommand struct {
	Meta
	modifyIndex	string
	dataFlags	string
}

func (c *KVWriteCommand) Help() string {
	helpText := `
Usage: consul-cli kv-write [options] path value

  Write a value to a given path.

Options:
` + c.ConsulHelp() + 
`  --modifyindex=<ModifyIndex>	Perform a Check-and-Set write.
				(default: not set)
  --flags=<number>		Integer value between 0 and 2^64 - 1
				(default: not set)
`

	return strings.TrimSpace(helpText)
}

func (c *KVWriteCommand) Run(args []string) int {
	c.AddDataCenter()
	flags := c.Meta.FlagSet()
	flags.StringVar(&c.modifyIndex, "modifyindex", "", "")
	flags.StringVar(&c.dataFlags, "flags", "", "")
	flags.Usage = func() { c.UI.Output(c.Help()) }

	if err := flags.Parse(args); err != nil {
		return 1
	}

	extra := flags.Args()
	if len(extra) < 2 {
		c.UI.Error("Key path and value must be specified")
		c.UI.Error("")
		c.UI.Error(c.Help())
		return 1
	}

	path := extra[0]
	value := strings.Join(extra[1:], " ")

	kv := new(consulapi.KVPair)

	kv.Key = path
	if strings.HasPrefix(value, "@") {
		v, err := ioutil.ReadFile(value[1:])
		if err != nil {
			c.UI.Error(fmt.Sprintf("ReadFile error: %v", err))
			return 1
		}
		kv.Value = v
	} else {
		kv.Value = []byte(value)
	}

	// &flags=
	//
	if c.dataFlags != "" {
		f, err := strconv.ParseUint(c.dataFlags, 0, 64)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing flags: %v", c.dataFlags))
			c.UI.Error("")
			return 1
		}
		kv.Flags = f
	}

	consul, err := c.Client()
	if err != nil {	
		c.UI.Error(err.Error())
		return 1
	}
	client := consul.KV()

	writeOpts := c.WriteOptions()

	if c.modifyIndex == "" {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(c.modifyIndex, 0, 64)
		if err != nil {
			c.UI.Error(fmt.Sprintf("Error parsing modifyIndex: %v", c.modifyIndex))
			c.UI.Error("")
			return 1
		}
		kv.ModifyIndex = i

		success, _, err := client.CAS(kv, writeOpts)
		if err != nil {
			c.UI.Error(err.Error())
			return 1
		}

		if !success {
			return 1
		}
	}


	return 0
}

func (c *KVWriteCommand) Synopsis() string {
	return "Write a value"
}
