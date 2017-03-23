package action

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl"
)

type raw struct {
	data string
}

func (r *raw) isSet() bool {
	return r.data != ""
}

// Return string of file contents
//
func (r *raw) readString() (string, error) {
	data, err := r.read()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Unmarshal JSON file contents
//
func (r *raw) readJSON(v interface{}) error {
	data, err := r.read()
	if err != nil {
		return err
	}

	// Try JSON first
	if err := json.Unmarshal(data, v); err == nil {
		return nil
	}

	// HCL next
	if err := hcl.Unmarshal(data, v); err == nil {
		return nil
	}

	return fmt.Errorf("Unable to unmarshal raw file")
}

// Read the data from a given path. Read from stdin
// if path == "-"
//
func (r *raw) read() ([]byte, error) {
	if r.data == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data, err := ioutil.ReadFile(r.data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
