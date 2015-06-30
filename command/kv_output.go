package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	consulapi "github.com/hashicorp/consul/api"
)

// Output structure
type KVOutput struct {
	ui		cli.Ui
	All		bool
	Key		bool
	CreateIndex	bool
	ModifyIndex	bool
	LockIndex	bool
	Flags		bool
	Value		bool
	Session		bool
}

// Output format structure
//
type OutputFormat struct {
	Type		string
	Delimiter	string
	Header		bool
}

// Conveninece structure for JSON
type KVJson struct {
	Key		string
	CreateIndex	uint64
	ModifyIndex	uint64
	LockIndex	uint64
	Flags		uint64
	Value		string
	Session		string
}

func NewKVOutput(ui cli.Ui, fields string) *KVOutput {
	kvo := new(KVOutput)

	kvo.ui = ui

	for _, field := range strings.Split(fields, ",") {
		f := strings.ToLower(field)

		switch {
		case f == "all":
			kvo.All = true
		case f == "key":
			kvo.Key = true
		case f == "createindex":
			kvo.CreateIndex = true
		case f == "modifyindex":
			kvo.ModifyIndex = true
		case f == "lockindex":
			kvo.LockIndex = true
		case f == "flags":
			kvo.Flags = true
		case f == "value":
			kvo.Value = true
		case f == "session":
			kvo.Session = true
		default:
			kvo.ui.Warn(fmt.Sprintf("Ignoring invalid field: %s", field))
		}
	}

	return kvo
}

func (kvo *KVOutput) Output(kv *consulapi.KVPair, of OutputFormat) int {
	f := strings.ToLower(of.Type)
	switch {
	case f == "json":
		return kvo.OutputJSON(kv, false)
	case f == "prettyjson":
		return kvo.OutputJSON(kv, true)
	case f == "text":
		return kvo.OutputText(kv, of)
	default: 
		kvo.ui.Error(fmt.Sprintf("Invalid output format: '%s'", of.Type))
	}

	return 1
}

func (kvo *KVOutput) OutputList(kvs *consulapi.KVPairs, of OutputFormat) int {
	f := strings.ToLower(of.Type)
	switch {
	case f == "json":
		return kvo.OutputJSONList(kvs, false)
	case f == "prettyjson":
		return kvo.OutputJSONList(kvs, true)
	case f == "text":
		return kvo.OutputTextList(kvs, of)
	default: 
		kvo.ui.Error(fmt.Sprintf("Invalid output format: '%s'", of.Type))
	}

	return 1
}

func (kvo *KVOutput) OutputJSONList(kvs *consulapi.KVPairs, prettyFlag bool) int {
	var err error
	var jsonRaw []byte

	kvjs := make([]*KVJson, len(*kvs))
	for i, kv := range *kvs {
		kvjs[i] = convertJSON(kv)
	}

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(kvjs, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(kvjs)
	}

	if err != nil {
		kvo.ui.Error(err.Error())
		return 1
	}

	kvo.ui.Output(string(jsonRaw))

	return 0
}

func (kvo *KVOutput) OutputJSON(kv *consulapi.KVPair, prettyFlag bool) int {
	var err error
	var jsonRaw []byte

	kvj := convertJSON(kv)

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(kvj, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(kvj)
	}

	if err != nil {
		kvo.ui.Error(err.Error())
		return 1
	}

	kvo.ui.Output(string(jsonRaw))

	return 0
}

func convertJSON(kv *consulapi.KVPair) *KVJson {
	return &KVJson{
		Key:		kv.Key,
		CreateIndex:	kv.CreateIndex,
		ModifyIndex:	kv.ModifyIndex,
		LockIndex:	kv.LockIndex,
		Flags:		kv.Flags,
		Value:		string(kv.Value),
		Session:	kv.Session,
	}
}

func (kvo *KVOutput) OutputText(kv *consulapi.KVPair, of OutputFormat) int {
	s := kvo.makeTextArray(kv)

	if of.Header {
		kvo.OutputHeader(of)
	}

	kvo.ui.Output(strings.Join(s, of.Delimiter))

	return 0
}

func (kvo *KVOutput) OutputTextList(kvs *consulapi.KVPairs, of OutputFormat) int {
	if of.Header {
		kvo.OutputHeader(of)
	}

	for _, kv := range *kvs {
		s := kvo.makeTextArray(kv)
		kvo.ui.Output(strings.Join(s, of.Delimiter))
	}

	return 0
}

func (kvo *KVOutput) makeTextArray(kv *consulapi.KVPair) []string {
	s := []string{}
	if kvo.Key || kvo.All {
		s = append(s, kv.Key)
	}
	if kvo.CreateIndex || kvo.All {
		s = append(s, fmt.Sprintf("%d", kv.CreateIndex))
	}
	if kvo.ModifyIndex || kvo.All {
		s = append(s, fmt.Sprintf("%d", kv.ModifyIndex))
	}
	if kvo.LockIndex || kvo.All {
		s = append(s, fmt.Sprintf("%d", kv.LockIndex))
	}
	if kvo.Flags || kvo.All {
		s = append(s, fmt.Sprintf("%d", kv.Flags))
	}
	if kvo.Value || kvo.All {
		s = append(s, string(kv.Value))
	}
	if kvo.Session || kvo.All {
		s = append(s, kv.Session)
	}

	return s
}


func (kvo *KVOutput) OutputHeader(of OutputFormat) {
	s := []string{}

	if kvo.Key || kvo.All {
		s = append(s, "key")
	}
	if kvo.CreateIndex || kvo.All {
		s = append(s, "createindex")
	}
	if kvo.ModifyIndex || kvo.All {
		s = append(s, "modifyindex")
	}
	if kvo.LockIndex || kvo.All {
		s = append(s, "lockindex")
	}
	if kvo.Flags || kvo.All {
		s = append(s, "flags")
	}
	if kvo.Value || kvo.All {
		s = append(s, "value")
	}
	if kvo.Session || kvo.All {
		s = append(s, "session")
	}

	kvo.ui.Output(fmt.Sprintf("#%s", strings.Join(s, of.Delimiter)))
}
