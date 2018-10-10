package action

import (
	"fmt"
	"strconv"
)

type mapSliceValue struct {
	value   *[]map[string]interface{}
	current map[string]interface{}
}

func newMapSliceValue(p *[]map[string]interface{}) *mapSliceValue {
	msv := new(mapSliceValue)
	msv.value = p
	return msv
}

func (msv *mapSliceValue) Set(val string) error {
	n := make(map[string]interface{})
	*msv.value = append(*msv.value, n)
	msv.current = n
	return nil
}

func (msv *mapSliceValue) Type() string {
	return "mapSlice"
}

func (msv *mapSliceValue) String() string {
	return fmt.Sprintf("%v", msv.value)
}

func (msv *mapSliceValue) IsBoolFlag() bool { return true }

type mapValue struct {
	msv    *mapSliceValue
	member string
	t      string
}

func newMapValue(p *mapSliceValue, member string, t string) *mapValue {
	mv := new(mapValue)
	mv.msv = p
	mv.member = member
	mv.t = t
	return mv
}

func (mv *mapValue) Set(val string) error {
	if mv.msv.current == nil {
		return fmt.Errorf("Trying to write flag to empty map")
	}

	switch mv.t {
	case "string":
		mv.msv.current[mv.member] = val
	case "stringSlice":
		if v, ok := mv.msv.current[mv.member]; ok {
			strings, _ := readAsCSV(v.(string))
			strings = append(strings, val)
			mv.msv.current[mv.member], _ = writeAsCSV(strings)
		} else {
			mv.msv.current[mv.member], _ = writeAsCSV([]string{val})
		}
	case "bool":
		v, err := strconv.ParseBool(val)
		mv.msv.current[mv.member] = v
		return err
	case "uint64":
		v, err := strconv.ParseUint(val, 0, 64)
		mv.msv.current[mv.member] = v
		return err
	}

	return nil
}

func (mv *mapValue) Type() string {
	return mv.t
}

func (mv *mapValue) String() string {
	return fmt.Sprintf("%s", mv.member)
}

func (mv *mapValue) IsBoolFlag() bool {
	return mv.t == "bool"
}
