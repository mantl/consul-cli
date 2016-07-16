package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

type Acl struct {
	*Cmd
}

func (root *Cmd) initAcl() {
	a := Acl{Cmd: root}

	aclCmd := &cobra.Command{
		Use:   "acl",
		Short: "Consul /acl endpoint interface",
		Long:  "Consul /acl endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	a.AddCloneSub(aclCmd)
	a.AddCreateSub(aclCmd)
	a.AddDestroySub(aclCmd)
	a.AddInfoSub(aclCmd)
	a.AddListSub(aclCmd)
	a.AddUpdateSub(aclCmd)

	a.AddCommand(aclCmd)
}

func (a *Acl) CheckIdArg(args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("ACL id must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one id allowed")
	}

	return nil
}

type ConfigRule struct {
	PathType string
	Path     string
	Policy   string
}

func (a *Acl) ParseRuleConfig(s string) (*ConfigRule, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, errors.New("cannot specify empty rule declaration")
	}

	var pathType, path, policy string
	parts := strings.Split(s, ":")

	switch len(parts) {
	case 2:
		pathType, path = parts[0], parts[1]
		policy = "read"
	case 3:
		pathType, path, policy = parts[0], parts[1], parts[2]
	default:
		return nil, fmt.Errorf("invalid rule declaration '%s'", s)
	}

	return &ConfigRule{pathType, path, policy}, nil
}

type rulePath struct {
	Policy string
}

type aclRule struct {
	Key     map[string]*rulePath `json:"key,omitempty"`
	Service map[string]*rulePath `json:"service,omitempty"`
	Event   map[string]*rulePath `json:"event,omitempty"`
	Query   map[string]*rulePath `json:"query,omitempty"`
}

func NewAclRule() *aclRule {
	return &aclRule{
		Key:     make(map[string]*rulePath),
		Service: make(map[string]*rulePath),
		Event:   make(map[string]*rulePath),
		Query:   make(map[string]*rulePath),
	}
}

// Convert a list of Rules to a JSON string
func (a *Acl) GetRulesString(rs []*ConfigRule) (string, error) {
	rules := NewAclRule()

	for _, r := range rs {
		// Verify policy is one of "read", "write", or "deny"
		policy := strings.ToLower(r.Policy)
		switch policy {
		case "read", "write", "deny":
		default:
			return "", fmt.Errorf("Invalid rule policy: '%s'", r.Policy)
		}

		switch strings.ToLower(r.PathType) {
		case "key":
			rules.Key[r.Path] = &rulePath{r.Policy}
		case "service":
			rules.Service[r.Path] = &rulePath{r.Policy}
		case "event":
			rules.Event[r.Path] = &rulePath{r.Policy}
		case "query":
			rules.Query[r.Path] = &rulePath{r.Policy}
		default:
			return "", fmt.Errorf("Invalid path type: '%s'", r.PathType)
		}
	}

	ruleBytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}

	return string(ruleBytes), nil
}
