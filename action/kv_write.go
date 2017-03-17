package action

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-multierror"
)

type kvWrite struct {
	modifyIndex string
	flags string

	*config
}

func KvWriteAction() Action {
	return &kvWrite{
		config: &gConfig,
	}
}

func (k *kvWrite) CommandFlags() *flag.FlagSet {
	f := newFlagSet()

	f.StringVar(&k.modifyIndex, "modifyindex", "", "Perform a Check-and-Set write")
	f.StringVar(&k.flags, "flags", "", "Integer value between 0 and 2^64 - 1")

	k.addDatacenterFlag(f)
	k.addRawFlag(f)

	return f
}

func (k *kvWrite) Run(args []string) error {
	// Handle raw write
	if k.raw.isSet() {
		return k.writeRaw()
	}

	if len(args) < 2 {
		return fmt.Errorf("Key path and value must be specified")
	}

	path := args[0]
	value := strings.Join(args[1:], " ")

	kv := new(consulapi.KVPair)

	kv.Key = path
	if strings.HasPrefix(value, "@") {
		v, err := ioutil.ReadFile(value[1:])
		if err != nil {
			return fmt.Errorf("ReadFile error: %v", err)
		}
		kv.Value = v
	} else {
		kv.Value = []byte(value)
	}

	// &flags=
	//
	if k.flags != "" {
		f, err := strconv.ParseUint(k.flags, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing flags: %v", k.flags)
		}
		kv.Flags = f
	}

	return k.write(kv)
}

func (k *kvWrite) writeRaw() error {
	data, err := k.raw.read()
	if err != nil {
		return err
	}

	var kvp consulapi.KVPair
	if err := json.Unmarshal(data, &kvp); err == nil {
		return k.write(&kvp)
	}

	var kvps consulapi.KVPairs
	if err := json.Unmarshal(data, &kvps); err == nil {
		var result error
		for _, kvp := range kvps {
			if err := k.write(kvp); err != nil {
				result = multierror.Append(result, err)
			}
		}

		return result
	}

	return fmt.Errorf("Unable to unmarshal raw data to KVPair or KVPairs")
}

func (k *kvWrite) write(kv *consulapi.KVPair) error {
	client, err := k.newKv()
	if err != nil {
		return err
	}
	writeOpts := k.writeOptions()

	if k.modifyIndex == "" {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(k.modifyIndex, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing modifyIndex: %v", k.modifyIndex)
		}
		kv.ModifyIndex = i

		success, _, err := client.CAS(kv, writeOpts)
		if err != nil {
			return err
		}

		if !success {
			return fmt.Errorf("Failed to write to K/V")
		}
	}

	return nil
}


