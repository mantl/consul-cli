package commands

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getStringSlice works around viper.GetStringSlice's bad behavior.
// viper.GetStringSlice gets the flag with viper.Get() and casts the
// result use spf13.cast.ToStringSlice(). The issue is that viper.Get()
// returns an interface{} which cast.ToStringSlice() converts to
// []string{value}
//
// so, for example, `--arg=foo --arg=bar` using viper.GetStringSlice()
// will return []string{"foo,bar"} which is clearly incorrect.
//
func getStringSlice(cmd *cobra.Command, flag string) []string {
	ss, err := cmd.Flags().GetStringSlice(flag)
	if err != nil {
		return []string{}
	}

	return ss
}

// viper does have a GetUint64 function. Create one here.
func getUint64(key string) uint64 {
	return toUint64(viper.Get(key))
}

func toUint64(i interface{}) uint64 {
	v, _ := toUint64E(i)
	return v
}

// Modeled after spf13/cast ToInt64E
func toUint64E(i interface{}) (uint64, error) {
	i = indirect(i)

	switch s := i.(type) {
	case uint64:
		return s, nil
	case uint:
		return uint64(s), nil
	case uint32:
		return uint64(s), nil
	case uint16:
		return uint64(s), nil
	case uint8:
		return uint64(s), nil
	case string:
		v, err := strconv.ParseUint(s, 0, 0)
		if err == nil {
			return v, nil
		}

		return 0, fmt.Errorf("Unable to cast %#v to uint64", i)
	case float64:
		return uint64(s), nil
	case bool:
		if bool(s) {
			return uint64(1), nil
		}
		return uint64(0), nil
	case nil:
		return uint64(0), nil
	default:
		return uint64(0), fmt.Errorf("Unable to cast %#v to uint64", i)
	}
}

// From html/template/content.go
// Copyright 2011 The Go Authors. All rights reserved.
// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}
