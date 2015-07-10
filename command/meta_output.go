package command

import (
	"encoding/json"
)

func (m *Meta) OutputJSON(v interface{}, prettyFlag bool) {
	var err error
	var jsonRaw []byte

	if prettyFlag {
		jsonRaw, err = json.MarshalIndent(v, "", "  ")
	} else {
		jsonRaw, err = json.Marshal(v)
	}

	if err != nil {
		m.UI.Output(err.Error())
		return
	}

	m.UI.Output(string(jsonRaw))
}

