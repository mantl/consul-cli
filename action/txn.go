package action

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type txn struct {
	kv []map[string]interface{}

	*config
}

func TxnAction() Action {
	return &txn{
		config: &gConfig,
	}
}

func (t *txn) CommandFlags() *flag.FlagSet {
	f := t.newFlagSet(FLAG_DATACENTER, FLAG_RAW)

	kmv := newMapSliceValue(&t.kv)

	f.Var(kmv, "kv", "Begin a KV txn")
	f.Var(newMapValue(kmv, "verb", "string"), "verb", "Type of operation to perform")
	f.Var(newMapValue(kmv, "key", "string"), "key", "Full path of the entry")
	f.Var(newMapValue(kmv, "value", "string"), "value", "Entry value. Use @filename to read from file")
	f.Var(newMapValue(kmv, "flags", "uint64"), "flags", "Integer value between 0 and 2^64 - 1")
	f.Var(newMapValue(kmv, "index", "uint64"), "index", "Modify index for CAS operations")
	f.Var(newMapValue(kmv, "session", "string"), "session", "Session ID for locking/unlocking")

	return f
}

func (t *txn) Run(args []string) error {
	var ops consulapi.KVTxnOps

	if t.raw.isSet() {
		// To match the documentation for the /v1/txn PUT operation, we read
		// the raw JSON as a slice of TxnOps. We then copy the KV ops to a slice
		// of KVTXnOp.
		var txnops consulapi.TxnOps

		if err := t.raw.readJSON(&txnops); err != nil {
			return err
		}

		ops = make([]*consulapi.KVTxnOp, len(txnops))

		for i, to := range txnops {
			ops[i] = to.KV
		}
	} else {
		if len(t.kv) == 0 {
			return fmt.Errorf("No KV operations specified")
		}

		ops = make([]*consulapi.KVTxnOp, len(t.kv))
		for i, k := range t.kv {
			kto := new(consulapi.KVTxnOp)

			verb, ok := k["verb"]
			if !ok {
				return fmt.Errorf("No verb specified for KV Op #%d", i)
			}

			switch strings.ToLower(verb.(string)) {
			case "set":
				kto.Verb = consulapi.KVSet
			case "delete":
				kto.Verb = consulapi.KVDelete
			case "delete-cas":
				kto.Verb = consulapi.KVDeleteCAS
			case "delete-tree":
				kto.Verb = consulapi.KVDeleteTree
			case "cas":
				kto.Verb = "KVCAS"
			case "lock":
				kto.Verb = "KVLock"
			case "unlock":
				kto.Verb = "KVUnlock"
			case "get":
				kto.Verb = "KVGet"
			case "get-tree":
				kto.Verb = "KVGetTree"
			case "check-session":
				kto.Verb = "KVCheckSession"
			case "check-index":
				kto.Verb = "KVCheckIndex"
			}

			if v, ok := k["key"]; ok {
				kto.Key = v.(string)
			}
			if v, ok := k["value"]; ok {
				kto.Value = []byte(v.(string))
			}
			if v, ok := k["flags"]; ok {
				kto.Flags = v.(uint64)
			}
			if v, ok := k["index"]; ok {
				kto.Index = v.(uint64)
			}
			if v, ok := k["session"]; ok {
				kto.Session = v.(string)
			}

			ops[i] = kto
		}
	}

	client, err := t.newKv()
	if err != nil {
		return err
	}

	queryOpts := t.queryOptions()

	ok, response, _, err := client.Txn(ops, queryOpts)
	if err != nil {
		return err
	}

	if !ok {
		fmt.Println("Not okay")
	}

	rbyte, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(rbyte))

	return nil
}
