package commands

import (
	"encoding/json"
	"fmt"
	"path"
	"reflect"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Bulkload functions

func newKvBulkloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulkload",
		Short: "Bulkload value to the K/V store",
		Long:  "Bulkload value to the K/V store",
		RunE:  kvBulkload,
	}

	cmd.Flags().String("json", "", "Path to a JSON file to import")
	cmd.Flags().String("prefix", "", "Base K/V prefix")
	addDatacenterOption(cmd)
	addRawOption(cmd)

	return cmd
}

func kvBulkload(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var kvPairs []*consulapi.KVPair
	var err error

	switch {
	case viper.GetString("json") != "":
		kvPairs, err = flatKVToPairs(viper.GetString("json"))
		if err != nil {
			return err
		}
	case viper.GetString("raw") != "":
		kvPairs, err = rawDataToPairs(viper.GetString("raw"))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No bulkload method specified")
	}

	client, err := newKv()
	if err != nil {
		return err
	}
	writeOpts := writeOptions()

	prefix := viper.GetString("prefix")
	for _, kv := range kvPairs {
		// Prepend the requested prefix
		if prefix != "" {
			kv.Key = path.Join(prefix, kv.Key)
		}

		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

// Convert { "key": "value" } formatted JSON values
//
func flatKVToPairs(path string) ([]*consulapi.KVPair, error) {
	rawData, err := readRawData(path)
	if err != nil {
		return nil, err
	}

	data := map[string]interface{}{}
	if err := json.Unmarshal(rawData, &data); err != nil {
		return nil, err
	}

	return flatten(data)
}

func flatten(src map[string]interface{}) ([]*consulapi.KVPair, error) {
	return _flatten(reflect.ValueOf(src), "")
}

func _flatten(v reflect.Value, key string) ([]*consulapi.KVPair, error) {
	dest := []*consulapi.KVPair{}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Kind() != reflect.String {
				continue
			}
			kvps, err := _flatten(v.MapIndex(k), path.Join(key, k.String()))
			if err != nil {
				return nil, err
			}
			dest = append(dest, kvps...)
		}
	case reflect.Slice:
		for k := 0; k < v.Len(); k++ {
			kvps, err := _flatten(v.Index(k), path.Join(key, fmt.Sprintf("%d", k)))
			if err != nil {
				return nil, err
			}
			dest = append(dest, kvps...)
		}
	case reflect.Bool, reflect.Float64, reflect.String:
		dest = append(dest, &consulapi.KVPair{
			Key:   key,
			Value: []byte(fmt.Sprintf("%v", v.Interface())),
		})
	case reflect.Invalid:
		// JSON null
		dest = append(dest, &consulapi.KVPair{
			Key:   key,
			Value: []byte(fmt.Sprintf("%v", nil)),
		})
	default:
		return nil, fmt.Errorf("invalid kind: %s\n", v.Kind().String())
	}

	return dest, nil
}

// Read consulapi.KVPair and consul.KVPairs formatted JSON files
// for bulkload.
//
func rawDataToPairs(path string) ([]*consulapi.KVPair, error) {
	data, err := readRawData(path)
	if err != nil {
		return nil, err
	}

	var kvp consulapi.KVPair
	if err := json.Unmarshal(data, &kvp); err == nil {
		return []*consulapi.KVPair{ &kvp }, nil
	}

	var kvps consulapi.KVPairs
	if err := json.Unmarshal(data, &kvps); err == nil {
		rval := []*consulapi.KVPair{}
		for _, kvp := range kvps {
			rval = append(rval, kvp)
		}
		return rval, nil
	}

	return nil, fmt.Errorf("Unable to unmarshal raw data to KVPair or KVPairs")
}
