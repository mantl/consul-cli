package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

// Output structure
type kvOutput struct {
	Out io.Writer
	Err io.Writer

	All         bool
	Key         bool
	CreateIndex bool
	ModifyIndex bool
	LockIndex   bool
	Flags       bool
	Value       bool
	Session     bool
}

// Output format structure
//
type outputFormat struct {
	Type      string
	Delimiter string
	Header    bool
}

// Conveninece structure for JSON
type kvJson struct {
	Key         string
	CreateIndex uint64
	ModifyIndex uint64
	LockIndex   uint64
	Flags       uint64
	Value       string
	Session     string
}

func newKvOutput(out, err io.Writer, fields string) *kvOutput {
	kvo := new(kvOutput)

	kvo.Out = out
	kvo.Err = err

	for _, field := range strings.Split(fields, ",") {
		f := strings.ToLower(field)

		switch {
		case f == "":
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
			fmt.Fprintf(kvo.Out, "Ignoring invalid field: %s\n", field)
		}
	}

	return kvo
}

func (kvo *kvOutput) output(kv *consulapi.KVPair, of outputFormat) error {
	f := strings.ToLower(of.Type)
	switch {
	case f == "json":
		return kvo.outputJSON(kv, false)
	case f == "prettyjson":
		return kvo.outputJSON(kv, true)
	case f == "text":
		return kvo.outputText(kv, of)
	default:
		fmt.Fprintf(kvo.Err, "Invalid output format: '%s'\n", of.Type)
	}

	return nil
}

func (kvo *kvOutput) outputList(kvs *consulapi.KVPairs, of outputFormat) error {
	f := strings.ToLower(of.Type)
	switch {
	case f == "json":
		return kvo.outputJSONList(kvs, false)
	case f == "prettyjson":
		return kvo.outputJSONList(kvs, true)
	case f == "text":
		return kvo.outputTextList(kvs, of)
	default:
		fmt.Fprintf(kvo.Err, "Invalid output format: '%s'\n", of.Type)
	}

	return nil
}

func (kvo *kvOutput) outputJSONList(kvs *consulapi.KVPairs, prettyFlag bool) error {
	var err error
	var jsonRaw []byte

	kvjs := make([]*kvJson, len(*kvs))
	for i, kv := range *kvs {
		kvjs[i] = convertJSON(kv)
	}

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(kvjs, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(kvjs)
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(kvo.Out, string(jsonRaw))

	return nil
}

func (kvo *kvOutput) outputJSON(kv *consulapi.KVPair, prettyFlag bool) error {
	var err error
	var jsonRaw []byte

	kvj := convertJSON(kv)

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(kvj, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(kvj)
	}

	if err != nil {
		return err
	}

	fmt.Fprintln(kvo.Out, string(jsonRaw))

	return nil
}

func convertJSON(kv *consulapi.KVPair) *kvJson {
	return &kvJson{
		Key:         kv.Key,
		CreateIndex: kv.CreateIndex,
		ModifyIndex: kv.ModifyIndex,
		LockIndex:   kv.LockIndex,
		Flags:       kv.Flags,
		Value:       string(kv.Value),
		Session:     kv.Session,
	}
}

func (kvo *kvOutput) outputText(kv *consulapi.KVPair, of outputFormat) error {
	s := kvo.makeTextArray(kv)

	if of.Header {
		kvo.outputHeader(of)
	}

	fmt.Fprintln(kvo.Out, strings.Join(s, of.Delimiter))

	return nil
}

func (kvo *kvOutput) outputTextList(kvs *consulapi.KVPairs, of outputFormat) error {
	if of.Header {
		kvo.outputHeader(of)
	}

	for _, kv := range *kvs {
		s := kvo.makeTextArray(kv)
		fmt.Fprintln(kvo.Out, strings.Join(s, of.Delimiter))
	}

	return nil
}

func (kvo *kvOutput) makeTextArray(kv *consulapi.KVPair) []string {
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

func (kvo *kvOutput) outputHeader(of outputFormat) {
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

	fmt.Fprintf(kvo.Out, "#%s\n", strings.Join(s, of.Delimiter))
}
