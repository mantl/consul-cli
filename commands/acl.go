package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

// Clone functions

func newAclCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone <token>",
		Short: "Create a new token from an existing one",
		Long:  "Create a new token from an existing one",
		RunE:  aclClone,
	}

	return cmd
}

func aclClone(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	newid, _, err := client.Clone(args[0], writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(newid)

	return nil
}

// Create functions

func newAclCreateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [<id>]",
		Short: "Create an ACL. Requires a management token.",
		Long:  "Create an ACL. Requires a management token.",
		RunE:  aclCreate,
	}

	cmd.Flags().Bool("management", false, "Create a management token")
	cmd.Flags().String("name", "", "Name of the ACL")
	cmd.Flags().StringSlice("rule", nil, "Rule to create. Can be multiple rules on a command line. Format is type:path:policy")

	addRawOption(cmd)

	return cmd
}

func aclCreate(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = viper.GetString("name")

	switch {
	case len(args) == 1:
		entry.ID = args[0]
	case len(args) > 1:
		return fmt.Errorf("Only one ACL identified can be specified")
	}

	if viper.GetBool("management") {
		entry.Type = consulapi.ACLManagementType
	} else {
		entry.Type = consulapi.ACLClientType
	}

	if raw := viper.GetString("raw"); raw != "" {
		rules, err := readRawString(raw)
		if err != nil {
			return err
		}
		entry.Rules = rules
	} else {
		rules, err := getRulesString(getStringSlice(cmd, "rule"))
		if err != nil {
			return err
		}
		entry.Rules = rules
	}

	writeOpts := writeOptions()
	id, _, err := client.Create(entry, writeOpts)
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

// Destroy functions

func newAclDestroyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy <token>",
		Short: "Destroy an ACL",
		Long:  "Destroy an ACL",
		RunE:  aclDestroy,
	}

	return cmd
}

func aclDestroy(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()
	_, err = client.Destroy(id, writeOpts)
	if err != nil {
		return err
	}

	return nil
}

// Info functions

func newAclInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <token>",
		Short: "Query information about an ACL token",
		Long:  "Query information about an ACL token",
		RunE:  aclInfo,
	}

	addTemplateOption(cmd)

	return cmd
}

func aclInfo(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	acl, _, err := client.Info(id, queryOpts)
	if err != nil {
		return err
	}

	return output(acl)
}

// List functions

func newAclListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all active ACL tokens",
		Long:  "List all active ACL tokens",
		RunE:  aclList,
	}

	addTemplateOption(cmd)

	return cmd
}

func aclList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	acls, _, err := client.List(queryOpts)
	if err != nil {
		return err
	}

	return output(acls)
}

// Replication functions

func newAclReplicationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "replication",
		Short:  "Get the status of the ACL replication process",
		Long:   "Get the status of the ACL replication process",
		RunE:   aclReplication,
		Hidden: true,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func aclReplication(cmd *cobra.Command, args []string) error {
	return fmt.Errorf("ACL replication status not available in Consul API")

	//	viper.BindPFlags(cmd.Flags())
	//
	//	client, err := newACL()
	//	if err != nil {
	//		return err
	//	}
}

// Update functions

func newAclUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <token>",
		Short: "Update an ACL. Will be created if it doesn't exist",
		Long:  "Update an ACL. Will be created if it doesn't exist",
		RunE:  aclUpdate,
	}

	cmd.Flags().Bool("management", false, "Create a management token")
	cmd.Flags().String("name", "", "Name of the ACL")
	cmd.Flags().StringSlice("rule", nil, "Rule to update. Can be multiple rules on a command line. Format is type:path:policy")

	addRawOption(cmd)

	return cmd
}

func aclUpdate(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("An ACL id must be specified")
	}
	id := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newACL()
	if err != nil {
		return err
	}

	entry := new(consulapi.ACLEntry)
	entry.Name = viper.GetString("name")
	entry.ID = id

	if viper.GetBool("management") {
		entry.Type = consulapi.ACLManagementType
	} else {
		entry.Type = consulapi.ACLClientType
	}

	if raw := viper.GetString("raw"); raw != "" {
		rules, err := readRawString(raw)
		if err != nil {
			return err
		}
		entry.Rules = rules
	} else {
		rules, err := getRulesString(getStringSlice(cmd, "rule"))
		if err != nil {
			return err
		}
		entry.Rules = rules
	}

	writeOpts := writeOptions()
	_, err = client.Update(entry, writeOpts)
	if err != nil {
		return err
	}

	return nil
}
