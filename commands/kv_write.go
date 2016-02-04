package commands

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	consulapi "github.com/hashicorp/consul/api"
)

type KvWriteOptions struct {
	ModifyIndex	string
	DataFlags	string
}

func (k *Kv) AddWriteSub(cmd *cobra.Command) {
	kwo := &KvWriteOptions{}

	writeCmd := &cobra.Command{
		Use: "write <path> <value>",
		Short: "Write a value to a given path",
		Long: "Write a value to a given path",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Write(args, kwo)
		},
	}

	oldWriteCmd := &cobra.Command{
		Use: "kv-write <path> <value>",
		Short: "Write a value to a given path",
		Long: "Write a value to a given path",
		Deprecated: "Use kv write",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Write(args, kwo)
		},
	}

	writeCmd.Flags().StringVar(&kwo.ModifyIndex, "modifyindex", "", "Perform a Check-and-Set write")
	writeCmd.Flags().StringVar(&kwo.DataFlags, "flags", "", "Integer value between 0 and 2^64 - 1")
	k.AddDatacenterOption(writeCmd)

	oldWriteCmd.Flags().StringVar(&kwo.ModifyIndex, "modifyindex", "", "Perform a Check-and-Set write")
	oldWriteCmd.Flags().StringVar(&kwo.DataFlags, "flags", "", "Integer value between 0 and 2^64 - 1")
	k.AddDatacenterOption(oldWriteCmd)

	cmd.AddCommand(writeCmd)

	k.AddCommand(oldWriteCmd)
}

func (k *Kv) Write(args []string, kwo *KvWriteOptions) error {
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
	if kwo.DataFlags != "" {
		f, err := strconv.ParseUint(kwo.DataFlags, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing flags: %v", kwo.DataFlags)
		}
		kv.Flags = f
	}

	consul, err := k.Client()
	if err != nil {	
		return err
	}
	client := consul.KV()

	writeOpts := k.WriteOptions()

	if kwo.ModifyIndex == "" {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(kwo.ModifyIndex, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing modifyIndex: %v", kwo.ModifyIndex)
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
