package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

func (c *Cmd) Output(v interface{}) error {
	if c.Template == "" {
		return c.OutputJSON(v, true)
	} else {
		return c.OutputTemplate(v)
	}

	return nil
}

func (c *Cmd) OutputJSON(v interface{}, prettyFlag bool) error {
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

	fmt.Fprintf(c.Out, string(jsonRaw))

	return nil
}

func (c *Cmd) OutputTemplate(v interface{}) error {
	if c.Template == "" {
		return fmt.Errorf("Empty output template")
	}

	if strings.HasPrefix(c.Template, "@") {
		v, err := ioutil.ReadFile(c.Template[1:])
		if err != nil {
			return err
		}
		c.Template = string(v)
	}

	template, err := template.New("").Parse(c.Template)
	if err != nil {
		return err
	}

	return template.Execute(c.Out, v)
}
