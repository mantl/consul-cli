package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	consulapi "github.com/hashicorp/consul/api"
)

type output struct {
	template string

	fields    string
	format    string
	delimiter string
	header    bool
}

func (o *output) output(v interface{}) error {
	if o.template == "" {
		return o.outputJSON(v, true)
	} else {
		return o.outputTemplate(v)
	}
}

func (o *output) outputKv(v interface{}) error {
	if o.template == "" {
		switch strings.ToLower(o.format) {
		case "json":
			return o.outputJSON(v, false)
		case "prettyjson":
			return o.outputJSON(v, true)
		case "text":
			return o.outputKvText(v)
		default:
			return fmt.Errorf("Invalid output format: '%s'\n", o.format)
		}
	} else {
		return o.outputTemplate(v)
	}
}

func (o *output) outputJSON(v interface{}, prettyFlag bool) error {
	var err error
	var jsonRaw []byte

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(v, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(v)
	}

	if err != nil {
		return err
	}

	jsonStr := string(jsonRaw)
	if strings.HasSuffix(jsonStr, "\n") {
		fmt.Printf(jsonStr)
	} else {
		fmt.Printf(jsonStr + "\n")
	}

	return nil
}

func (o *output) outputTemplate(v interface{}) error {
	tmpl := o.template
	if tmpl == "" {
		return fmt.Errorf("Empty output template")
	}

	if strings.HasPrefix(tmpl, "@") {
		v, err := ioutil.ReadFile(tmpl[1:])
		if err != nil {
			return err
		}
		tmpl = string(v)
	}

	template, err := template.New("").Parse(tmpl)
	if err != nil {
		return err
	}

	return template.Execute(os.Stdout, v)
}

// KV text output functions

var validFields = map[string]string{
	"createindex": "CreateIndex",
	"flags":       "Flags",
	"key":         "Key",
	"lockindex":   "LockIndex",
	"modifyindex": "ModifyIndex",
	"session":     "Session",
	"value":       "Value",
}

// Add the header line to the output template if requested
func (o *output) addHeader(b *bytes.Buffer) error {
	if !o.header {
		return nil
	}

	fieldList := o.fields

	// All field case
	if fieldList == "all" || fieldList == "" {
		fieldList = "key,value,flags,createindex,modifyindex,lockindex,session"
	}

	delimiter := o.delimiter

	fields := strings.Split(fieldList, ",")
	for i, field := range fields {
		if i >= 1 {
			if _, err := b.WriteString(delimiter); err != nil {
				return err
			}
		}

		if f, ok := validFields[strings.ToLower(field)]; !ok {
			return fmt.Errorf("Invalid field: %s", field)
		} else {
			b.WriteString(f)
		}
	}
	b.WriteString("\n")

	return nil
}

func (o *output) addFields(b *bytes.Buffer) error {
	fieldList := o.fields

	// All field case
	if fieldList == "all" || fieldList == "" {
		fieldList = "key,value,flags,createindex,modifyindex,lockindex,session"
	}

	delimiter := o.delimiter
	for i, field := range strings.Split(fieldList, ",") {
		if i >= 1 {
			if _, err := b.WriteString(delimiter); err != nil {
				return err
			}
		}

		if f, ok := validFields[strings.ToLower(field)]; !ok {
			return fmt.Errorf("Invalid field: %s", field)
		} else {
			// Special []byte handling
			if field == "value" {
				b.WriteString("{{printf \"%s\" ." + f + "}}")
			} else {
				b.WriteString("{{." + f + "}}")
			}
		}
	}
	b.WriteString("\n")

	return nil
}

func (o *output) outputKvText(v interface{}) error {
	b := new(bytes.Buffer)

	var isList bool

	switch v.(type) {
	case *consulapi.KVPairs:
		isList = true
	case *consulapi.KVPair:
		isList = false
	default:
		return fmt.Errorf("Invalid KVPair type: %t\n", v)
	}

	if err := o.addHeader(b); err != nil {
		return err
	}

	if isList {
		b.WriteString("{{range .}}")
	}

	if err := o.addFields(b); err != nil {
		return err
	}

	if isList {
		b.WriteString("{{end}}")
	}

	o.template = b.String()
	return o.outputTemplate(v)
}
