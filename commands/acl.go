package commands

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/spf13/cobra"
)

func newAclCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "Consul /acl endpoint interface",
		Long:  "Consul /acl endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newAclCloneCommand())
	cmd.AddCommand(newAclCreateCommand())
	cmd.AddCommand(newAclDestroyCommand())
	cmd.AddCommand(newAclInfoCommand())
	cmd.AddCommand(newAclListCommand())
	cmd.AddCommand(newAclReplicationCommand())
	cmd.AddCommand(newAclUpdateCommand())

	return cmd
}

type rulePath struct {
	Policy string
}

type aclRule struct {
	Key      map[string]*rulePath `json:"key,omitempty"`
	Service  map[string]*rulePath `json:"service,omitempty"`
	Event    map[string]*rulePath `json:"event,omitempty"`
	Query    map[string]*rulePath `json:"query,omitempty"`
	Keyring  string               `json:"keyring,omitempty"`
	Operator string               `json:"operator,omitempty"`
}

// getPolicy return "read" if the index i is not set in the
// rs array.
func getPolicy(rs []string, i int) string {
	if i >= len(rs) {
		return "read"
	}

	return rs[i]
}

// getPath returns "" if the inde i is not set in the rs array
func getPath(rs []string, i int) string {
	if i >= len(rs) {
		return ""
	}

	return rs[i]
}

// Convert a list of Rules to a JSON string
func getRulesString(rs []string) (string, error) {
	if len(rs) <= 0 {
		return "", errors.New("No ACL rules specified")
	}

	rules := &aclRule{
		Key:     make(map[string]*rulePath),
		Service: make(map[string]*rulePath),
		Event:   make(map[string]*rulePath),
		Query:   make(map[string]*rulePath),
	}

	for _, r := range rs {
		if len(strings.TrimSpace(r)) < 1 {
			return "", errors.New("cannot specify empty rule declaration")
		}

		parts := strings.Split(r, ":")
		switch strings.ToLower(parts[0]) {
		case "operator":
			rules.Operator = getPolicy(parts, 1)
		case "keyring":
			rules.Keyring = getPolicy(parts, 1)
		case "key":
			rules.Key[getPath(parts, 1)] = &rulePath{Policy: getPolicy(parts, 2)}
		case "service":
			rules.Service[getPath(parts, 1)] = &rulePath{Policy: getPolicy(parts, 2)}
		case "event":
			rules.Event[getPath(parts, 1)] = &rulePath{Policy: getPolicy(parts, 2)}
		case "query":
			rules.Query[getPath(parts, 1)] = &rulePath{Policy: getPolicy(parts, 2)}
		}
	}

	ruleBytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}

	return string(ruleBytes), nil
}
