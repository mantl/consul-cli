package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Write functions

func newKvWriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write <path> <value>",
		Short: "Write a value to a given path",
		Long:  "Write a value to a given path",
		RunE:  kvWrite,
	}

	cmd.Flags().String("modifyindex", "", "Perform a Check-and-Set write")
	cmd.Flags().String("flags", "", "Integer value between 0 and 2^64 - 1")
	addDatacenterOption(cmd)
	addRawOption(cmd)

	return cmd
}

func kvWrite(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	// Handle raw write
	if raw := viper.GetString("raw"); raw != "" {
		return kvWriteRaw(raw)
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
	if flags := viper.GetString("flags"); flags != "" {
		f, err := strconv.ParseUint(flags, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing flags: %v", flags)
		}
		kv.Flags = f
	}

	return kvDoWrite(kv)
}

func kvWriteRaw(raw string) error {
	data, err := readRawData(raw)
	if err != nil {
		return err
	}

	var kvp consulapi.KVPair
	if err := json.Unmarshal(data, &kvp); err == nil {
		return kvDoWrite(&kvp)
	}

	var kvps consulapi.KVPairs
	if err := json.Unmarshal(data, &kvps); err == nil {
		var result error
		for _, kvp := range kvps {
			if err := kvDoWrite(kvp); err != nil {
				result = multierror.Append(result, err)
			}
		}

		return result
	}

	return fmt.Errorf("Unable to unmarshal raw data to KVPair or KVPairs")
}

func kvDoWrite(kv *consulapi.KVPair) error {
	client, err := newKv()
	if err != nil {
		return err
	}
	writeOpts := writeOptions()

	if viper.GetString("modifyindex") == "" {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(viper.GetString("modifyindex"), 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing modifyIndex: %v", viper.GetString("modifyindex"))
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


