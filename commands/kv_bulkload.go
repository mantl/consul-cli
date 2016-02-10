package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type KvBulkloadOptions struct {
	Json string
	Prefix string
}

func (k *Kv) AddBulkloadSub(cmd *cobra.Command) {
	kbo := &KvBulkloadOptions{}

	bulkloadCmd := &cobra.Command{
		Use:   "bulkload",
		Short: "Bulkload value to the K/V store",
		Long: "Bulkload value to the K/V store",
		RunE: func(cmd *cobra.Command, args []string) error {
			return k.Bulkload(args, kbo)
		},
	}

	bulkloadCmd.Flags().StringVar(&kbo.Json, "json", "", "Path to a JSON file to import")
	bulkloadCmd.Flags().StringVar(&kbo.Prefix, "prefix", "", "Base K/V prefix")
	k.AddDatacenterOption(bulkloadCmd)

	cmd.AddCommand(bulkloadCmd)
}

func (k *Kv) Bulkload(args []string, kbo *KvBulkloadOptions) error {
	var kvPairs []*consulapi.KVPair

	switch {
	case kbo.Json != "":
		rawData, err := ioutil.ReadFile(kbo.Json)
		if err != nil {
			return err
		}

		data := map[string]interface{}{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			return err
		}

		kvPairs = flatten(data, kbo.Prefix)
	default :
		return fmt.Errorf("Must specify --json")
	}
		
	client, err := k.KV()
	if err != nil {
		return err
	}
	writeOpts := k.WriteOptions()

	for _, kv := range kvPairs {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func flatten(src map[string]interface{}, prefix string) []*consulapi.KVPair {
	return _flatten(reflect.ValueOf(src), prefix)
}

func _flatten(v reflect.Value, prefix string) []*consulapi.KVPair {
	dest := []*consulapi.KVPair{}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
        case reflect.Map:
                for _, k := range v.MapKeys() {
                        if k.Kind() != reflect.String { continue }
                        dest = append(dest, _flatten(v.MapIndex(k), path.Join(prefix, k.String()))...)
                }
        case reflect.Slice:
                for k := 0; k < v.Len(); k++ {
                        dest = append(dest, _flatten(v.Index(k), path.Join(prefix, fmt.Sprintf("%d", k)))...)
                }
        case reflect.Bool, reflect.Float64, reflect.String:
                dest = append(dest, &consulapi.KVPair{
			Key: prefix,
			Value: []byte(fmt.Sprintf("%v", v.Interface())),
		})
        case reflect.Invalid:
                // JSON null
                dest = append(dest, &consulapi.KVPair{
			Key: prefix, 
			Value: []byte(fmt.Sprintf("%v", nil)),
		})
        default:
                fmt.Printf("invalid kind: %s\n", v.Kind().String())
        }

	return dest
}
