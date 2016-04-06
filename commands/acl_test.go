package commands

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestAcl() *Acl {
	return &Acl{}
}

func Test_ParseRuleConfig(t *testing.T) {
	cases := []struct{
		Value string
		ExpectedRule *ConfigRule
		ExpectedErr error
	}{
		{
			"key::write",
			&ConfigRule{"key","","write"},
			nil,
		},
		{
			"service:",
			&ConfigRule{"service","","read"},
			nil,
		},
		{
			"query:foo-:deny",
			&ConfigRule{"query","foo-","deny"},
			nil,
		},
		{
			"event:destroy-:write",
			&ConfigRule{"event","destroy-","write"},
			nil,
		},
		{
			"foo",
			nil,
			errors.New("error"),
		},
		{
			"foo:bar:baz:quux",
			nil,
			errors.New("error"),
		},
	}

	acl := newTestAcl()

	for _, c := range cases {
		ActualRule, ActualErr := acl.ParseRuleConfig(c.Value)
		if c.ExpectedRule == nil {
			assert.Nil(t, ActualRule)
		} else {
			assert.Equal(t, c.ExpectedRule.PathType, ActualRule.PathType)
			assert.Equal(t, c.ExpectedRule.Path, ActualRule.Path)
			assert.Equal(t, c.ExpectedRule.Policy, ActualRule.Policy)
		}

		if c.ExpectedErr != nil {
			assert.NotNil(t, ActualErr)
		} else {
			assert.Nil(t, ActualErr)
		}
	}
}

func Test_GetRulesString(t *testing.T) {
	cases := []struct{
		Value []*ConfigRule
		ExpectedString string
		ExpectedErr error
	}{
		{
			Value: []*ConfigRule{
				&ConfigRule{"key","","write"},
			},
			ExpectedString: `{"key":{"":{"Policy":"write"}}}`,
			ExpectedErr: nil,
		},
		{
			Value: []*ConfigRule{
				&ConfigRule{"service","","read"},
			},
			ExpectedString: `{"service":{"":{"Policy":"read"}}}`,
			ExpectedErr: nil,
		},
		{
			Value: []*ConfigRule{
				&ConfigRule{"query","foo-","deny"},
			},
			ExpectedString: `{"query":{"foo-":{"Policy":"deny"}}}`,
			ExpectedErr: nil,
		},
		{
			Value: []*ConfigRule{
				&ConfigRule{"event","destroy-","write"},
			},
			ExpectedString: `{"event":{"destroy-":{"Policy":"write"}}}`,
			ExpectedErr: nil,
		},
	}

	acl := newTestAcl()

	for _, c := range cases {
		ActualString, ActualErr := acl.GetRulesString(c.Value)
		assert.Equal(t, c.ExpectedString, ActualString)

		if c.ExpectedErr != nil {
			assert.NotNil(t, ActualErr)
		} else {
			assert.Nil(t, ActualErr)
		}
	}
}
