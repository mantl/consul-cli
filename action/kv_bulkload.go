package action

import (
	"encoding/json"
	"flag"
	"fmt"
	"path"
	"reflect"

	consulapi "github.com/hashicorp/consul/api"
)

type kvBulkload struct {
	json   string
	prefix string

	*config
}

func KvBulkloadAction() Action {
	return &kvBulkload{
		config: &gConfig,
	}
}

func (k *kvBulkload) CommandFlags() *flag.FlagSet {
	f := k.newFlagSet(FLAG_DATACENTER, FLAG_RAW)

	f.StringVar(&k.json, "json", "", "Path to a JSON file to import")
	f.StringVar(&k.prefix, "prefix", "", "Base K/V prefix")

	return f
}

func (k *kvBulkload) Run(args []string) error {
	var kvPairs []*consulapi.KVPair
	var err error

	switch {
	case k.json != "":
		kvPairs, err = k.flatKVToPairs()
		if err != nil {
			return err
		}
	case k.raw.isSet():
		kvPairs, err = k.rawDataToPairs()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("No bulkload method specified")
	}

	client, err := k.newKv()
	if err != nil {
		return err
	}
	writeOpts := k.writeOptions()

	for _, kv := range kvPairs {
		// Prepend the requested prefix
		if k.prefix != "" {
			kv.Key = path.Join(k.prefix, kv.Key)
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
func (k *kvBulkload) flatKVToPairs() ([]*consulapi.KVPair, error) {
	rawData, err := k.raw.read()
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
func (k *kvBulkload) rawDataToPairs() ([]*consulapi.KVPair, error) {
	data, err := k.raw.read()
	if err != nil {
		return nil, err
	}

	var kvp consulapi.KVPair
	if err := json.Unmarshal(data, &kvp); err == nil {
		return []*consulapi.KVPair{&kvp}, nil
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
