package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func output(v interface{}) error {
	if viper.GetString("template") == "" {
		return outputJSON(v, true)
	} else {
		return outputTemplate(v)
	}
}

func outputKv(v interface{}) error {
	switch t := strings.ToLower(viper.GetString("format")); t {
	case "json":
		return outputJSON(v, false)
	case "prettyjson":
		return outputJSON(v, true)
	case "text":
		return outputKvText(v)
	default:
		return fmt.Errorf("Invalid output format: '%s'\n", t)
	}

	return nil
}

func outputJSON(v interface{}, prettyFlag bool) error {
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

func outputTemplate(v interface{}) error {
	tmpl := viper.GetString("template")
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
        "flags": "Flags",
        "key": "Key",
        "lockindex": "LockIndex",
        "modifyindex": "ModifyIndex",
        "session": "Session",
        "value": "Value",
}

// Add the header line to the output template if requested
func addHeader(t *bytes.Buffer) error {
        if !viper.GetBool("header") {
                return nil
        }

        fieldList := viper.GetString("fields")

        // All field case
        if fieldList == "all" || fieldList == "" {
                fieldList = "key,value,flags,createindex,modifyindex,lockindex,session"
        }

        delimiter := viper.GetString("delimiter")

        fields := strings.Split(fieldList, ",")
        for i, field := range fields {
                if i >= 1 {
                        if _, err := t.WriteString(delimiter); err != nil {
                                return err
                        }
                }

                if f, ok := validFields[strings.ToLower(field)]; !ok {
                        return fmt.Errorf("Invalid field: %s", field)
                } else {
                        t.WriteString(f)
                }
        }
        t.WriteString("\n")

        return nil
}
func addFields(t *bytes.Buffer) error {
        fieldList := viper.GetString("fields")

        // All field case
        if fieldList == "all" || fieldList == "" {
                fieldList = "key,value,flags,createindex,modifyindex,lockindex,session"
        }

        delimiter := viper.GetString("delimiter")
        for i, field := range strings.Split(fieldList, ",") {
                if i >= 1 {
                        if _, err := t.WriteString(delimiter); err != nil {
                                return err
                        }
                }

                if f, ok := validFields[strings.ToLower(field)]; !ok {
                        return fmt.Errorf("Invalid field: %s", field)
                } else {
                        // Special []byte handling
                        if field == "value" {
                                t.WriteString("{{printf \"%s\" ." + f + "}}")
                        } else {
                                t.WriteString("{{." + f + "}}")
                        }
                }
        }
        t.WriteString("\n")

        return nil
}

func outputKvText(v interface{}) error {
        t := new(bytes.Buffer)

	var isList bool

	switch v.(type) {
	case *consulapi.KVPairs:
		isList = true
	case *consulapi.KVPair:
		isList = false
	default:
		return fmt.Errorf("Invalid KVPair type: %t\n", v)
	}

        if err := addHeader(t); err != nil {
                return err
        }

	if isList {
	        t.WriteString("{{range .}}")
	}

        if err := addFields(t); err != nil {
                return err
        }

	if isList {
	        t.WriteString("{{end}}")
	}

        viper.Set("template", t.String())
        return outputTemplate(v)
}
