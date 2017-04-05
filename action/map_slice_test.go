package action

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFlagSet(m *[]map[string]interface{}) *flag.FlagSet {
	f := flag.NewFlagSet("consul-cli-test", flag.ExitOnError)

	mv := newMapSliceValue(m)

	f.Var(mv, "key", "Test structure trigger flag")
	f.Var(newMapValue(mv, "stringField1", "string"), "string-field1", "String field 1")
	f.Var(newMapValue(mv, "stringField2", "string"), "string-field2", "String field 2")
	f.Var(newMapValue(mv, "boolField", "bool"), "bool-field", "Bool field")
	f.Var(newMapValue(mv, "uint64Field", "uint64"), "uint64-field", "uint64 field")

	return f
}

type testCase struct {
	args []string
	res []result
}

type result struct {
	index int
	field string
	val interface{}
}

func TestMapSlice(t *testing.T) {
	cases := []testCase{
		{
			[]string{"--key", "--string-field1", "testing"},
			[]result{{0, "stringField1", "testing"}},
		},
		{
			[]string{ 
				"--key", "--string-field1", "testing",
				"--key", "--string-field2", "another",
			},
			[]result{ 
				{0, "stringField1", "testing"},
				{1, "stringField2", "another"},
			},
		},
		{
			[]string{ "--key", "--bool-field" },
			[]result{{0, "boolField", true}},
		},
		{
			[]string{"--key", "--uint64-field", "123456789"},
			[]result{{0, "uint64Field", 123456789}},
		},
	}

	for _, c := range cases {
		m := []map[string]interface{}{}
		f := testFlagSet(&m)
		err := f.Parse(c.args)

		assert.Nil(t, err)
		for _, e := range c.res {
			assert.NotEmpty(t, m[e.index][e.field])
			assert.True(t, assert.ObjectsAreEqualValues(m[e.index][e.field], e.val))
		}
	}
}
