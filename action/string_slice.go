package action

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"
)

var _ = fmt.Fprint

// -- stringSlice Value
type stringSliceValue struct {
	value   *[]string
	changed bool
}

func newStringSliceValue(p *[]string) *stringSliceValue {
	ssv := new(stringSliceValue)
	ssv.value = p
	return ssv
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

func writeAsCSV(vals []string) (string, error) {
	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), fmt.Sprintln()), nil
}

func (s *stringSliceValue) Set(val string) error {
	v, err := readAsCSV(val)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = v
	} else {
		*s.value = append(*s.value, v...)
	}
	s.changed = true
	return nil
}

func (s *stringSliceValue) Type() string {
	return "stringSlice"
}

func (s *stringSliceValue) String() string {
	str, _ := writeAsCSV(*s.value)
	return "[" + str + "]"
}
