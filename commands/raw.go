package commands

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

func addRawOption(cmd *cobra.Command) {
	cmd.Flags().String("raw", "", "Raw JSON data for upload")
}

// Return string of file contents
//
func readRawString(path string) (string, error) {
	data, err := readRawData(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Unmarshal JSON file contents
//
func readRawJSON(path string, v interface{}) error {
	data, err := readRawData(path)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

// Read the data from a given path. Read from stdin
// if path == "-"
//
func readRawData(path string) ([]byte, error) {
	if path == "-" {
		data, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}
