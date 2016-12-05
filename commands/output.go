package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

func output(v interface{}) error {
	if viper.GetString("template") == "" {
		return outputJSON(v, true)
	} else {
		return outputTemplate(v)
	}
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
