package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newKvCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kv",
		Short: "Consul /kv endpoint interface",
		Long:  "Consul /kv endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newKvBulkloadCommand())
	cmd.AddCommand(newKvDeleteCommand())
	cmd.AddCommand(newKvKeysCommand())
	cmd.AddCommand(newKvLockCommand())
	cmd.AddCommand(newKvReadCommand())
	cmd.AddCommand(newKvUnlockCommand())
	cmd.AddCommand(newKvWatchCommand())
	cmd.AddCommand(newKvWriteCommand())

	return cmd
}

// Bulkload functions

func newKvBulkloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bulkload",
		Short: "Bulkload value to the K/V store",
		Long:  "Bulkload value to the K/V store",
		RunE:  kvBulkload,
	}

	cmd.Flags().String("json", "", "Path to a JSON file to import")
	cmd.Flags().String("prefix", "", "Base K/V prefix")
	addDatacenterOption(cmd)

	return cmd
}

func kvBulkload(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var kvPairs []*consulapi.KVPair

	switch {
	case viper.GetString("json") != "":
		rawData, err := ioutil.ReadFile(viper.GetString("json"))
		if err != nil {
			return err
		}

		data := map[string]interface{}{}
		if err := json.Unmarshal(rawData, &data); err != nil {
			return err
		}

		kvPairs = flatten(data, viper.GetString("prefix"))
	default:
		return fmt.Errorf("Must specify --json")
	}

	client, err := newKv()
	if err != nil {
		return err
	}
	writeOpts := writeOptions()

	for _, kv := range kvPairs {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

func flatten(src map[string]interface{}, prefix string) []*consulapi.KVPair {
	return _flatten(reflect.ValueOf(src), prefix)
}

func _flatten(v reflect.Value, prefix string) []*consulapi.KVPair {
	dest := []*consulapi.KVPair{}

	if v.Kind() == reflect.Interface {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			if k.Kind() != reflect.String {
				continue
			}
			dest = append(dest, _flatten(v.MapIndex(k), path.Join(prefix, k.String()))...)
		}
	case reflect.Slice:
		for k := 0; k < v.Len(); k++ {
			dest = append(dest, _flatten(v.Index(k), path.Join(prefix, fmt.Sprintf("%d", k)))...)
		}
	case reflect.Bool, reflect.Float64, reflect.String:
		dest = append(dest, &consulapi.KVPair{
			Key:   prefix,
			Value: []byte(fmt.Sprintf("%v", v.Interface())),
		})
	case reflect.Invalid:
		// JSON null
		dest = append(dest, &consulapi.KVPair{
			Key:   prefix,
			Value: []byte(fmt.Sprintf("%v", nil)),
		})
	default:
		fmt.Printf("invalid kind: %s\n", v.Kind().String())
	}

	return dest
}

// Delete functions

func newKvDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <path>",
		Short: "Delete a given path from the K/V",
		Long:  "Delete a given path from the K/V",
		RunE:  kvDelete,
	}

	cmd.Flags().String("modifyindex", "", "Perform a Check-and-Set delete")
	cmd.Flags().Bool("recurse", false, "Perform a recursive delete")
	addDatacenterOption(cmd)

	return cmd
}

func kvDelete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	switch {
	case viper.GetBool("recurse"):
		_, err := client.DeleteTree(path, writeOpts)
		if err != nil {
			return err
		}
	case viper.GetString("modifyindex") != "":
		m, err := strconv.ParseUint(viper.GetString("modifyindex"), 0, 64)
		if err != nil {
			return err
		}
		kv := consulapi.KVPair{
			Key:         path,
			ModifyIndex: m,
		}

		success, _, err := client.DeleteCAS(&kv, writeOpts)
		if err != nil {
			return err
		}

		if !success {
			return fmt.Errorf("Failed deleting")
		}
	default:
		_, err := client.Delete(path, writeOpts)
		if err != nil {
			return err
		}
	}

	return nil
}

// Keys functions

func newKvKeysCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys <path>",
		Short: "List K/V keys",
		Long:  "List K/V keys",
		RunE:  kvKeys,
	}

	cmd.Flags().String("separator", "", "List keys only up to a given separator")
	addDatacenterOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func kvKeys(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	data, _, err := client.Keys(path, viper.GetString("separator"), queryOpts)
	if err != nil {
		return err
	}

	viper.Set("template", kv_outputTemplate)

	return output(data)
}

var kv_outputTemplate = `{{range .}}{{.}}
{{end}}`

// Lock functions

func newKvLockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock <path>",
		Short: "Acquire a lock on a given path",
		Long:  "Acquire a lock on a given path",
		RunE:  kvLock,
	}
	cmd.Flags().String("behavior", "release", "Lock behavior. One of 'release' or 'delete'")
	cmd.Flags().String("ttl", "", "Lock time to live")
	cmd.Flags().Duration("lock-delay", 15*time.Second, "Lock delay")
	cmd.Flags().String("session", "", "Previously created session to use for lock")

	addDatacenterOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func kvLock(cmd *cobra.Command, args []string) error {
	var lockOpts *consulapi.KVPair

	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	// Work around a Consul API bug that ignores LockDelay == 0
	lockDelay := viper.GetDuration("lock-delay")
	if lockDelay == 0 {
		lockDelay = time.Nanosecond
	}

	client, err := newKv()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()
	queryOpts := queryOptions()
	queryOpts.WaitTime = 15 * time.Second

	sessionClient, err := newSession()
	if err != nil {
		return err
	}

	if viper.GetString("session") == "" {
		// Create the Consul session
		se, _, err := sessionClient.CreateNoChecks(&consulapi.SessionEntry{
			Name:      "Session for consul-cli",
			LockDelay: lockDelay,
			Behavior:  viper.GetString("behavior"),
			TTL:       viper.GetString("ttl"),
		}, writeOpts)
		if err != nil {
			return err
		}

		viper.Set("session", se)
		viper.Set("__clean_session", true)
	}

	session := viper.GetString("session")

	// Set the session to renew periodically
	sessionRenew := make(chan struct{})
	go sessionClient.RenewPeriodic(viper.GetString("ttl"), session, nil, sessionRenew)
	defer func() {
		close(sessionRenew)
		sessionRenew = nil
	}()

WAIT:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		destroySession(sessionClient)
		return err
	}

	locked := false

	if kv != nil && kv.Session == session {
		goto HELD
	}
	if kv != nil && kv.Session != "" {
		queryOpts.WaitIndex = meta.LastIndex
		goto WAIT
	}

	// Node doesn't already exist
	if kv == nil {
		lockOpts = &consulapi.KVPair{
			Key:     path,
			Session: session,
		}
	} else {
		lockOpts = &consulapi.KVPair{
			Key:     kv.Key,
			Flags:   kv.Flags,
			Value:   kv.Value,
			Session: session,
		}
	}

	// Try to acquire the lock
	locked, _, err = client.Acquire(lockOpts, nil)
	if err != nil {
		destroySession(sessionClient)
		return err
	}

	if !locked {
		select {
		case <-time.After(5 * time.Second):
			goto WAIT
		}
	}

HELD:
	fmt.Println(session)

	return nil
}

// Destroy the session on error. Only performed when
// __clean_session == true
//
func destroySession(s *consulapi.Session) error {
	if viper.GetBool("__clean_session") {
		session := viper.GetString("session")
		writeOpts := writeOptions()
		_, err := s.Destroy(session, writeOpts)
		if err != nil {
			return fmt.Errorf("Session not destroyed: %s", session)
		}
	}

	return nil
}

// Read functions

func newKvReadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <path>",
		Short: "Read a value from a given path",
		Long:  "Read a value from a given path",
		RunE:  kvRead,
	}

	cmd.Flags().String("fields", "value", "Comma separated list of fields to return")
	cmd.Flags().String("format", "text", "Output format. Supported options: text, json, prettyjson")
	cmd.Flags().String("delimiter", " ", "Output field delimiter")
	cmd.Flags().Bool("header", false, "Output a header row for text format")
	cmd.Flags().Bool("recurse", false, "Perform a recursive read")

	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func kvRead(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	kvo := NewKVOutput(os.Stdout, os.Stderr, viper.GetString("fields"))

	if viper.GetBool("recurse") {
		kvlist, _, err := client.List(path, queryOpts)
		if err != nil {
			return err
		}

		if kvlist == nil {
			return nil
		}

		if viper.GetString("template") != "" {
			return output(kvlist)
		} else {
			return kvo.OutputList(&kvlist, OutputFormat{
				Type:      viper.GetString("format"),
				Delimiter: viper.GetString("delimiter"),
				Header:    viper.GetBool("header"),
			})
		}
	} else {
		kv, _, err := client.Get(path, queryOpts)
		if err != nil {
			return err
		}

		if kv == nil {
			return nil
		}

		if viper.GetString("template") != "" {
			return output(kv)
		} else {
			return kvo.Output(kv, OutputFormat{
				Type:      viper.GetString("format"),
				Delimiter: viper.GetString("delimiter"),
				Header:    viper.GetBool("header"),
			})
		}
	}
}

// Unlock functions

func newKvUnlockCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unlock <path>",
		Short: "Release a lock on a given path",
		Long:  "Release a lock on a given path",
		RunE:  kvUnlock,
	}

	cmd.Flags().String("session", "", "Session ID of the lock holder. Required")
	cmd.Flags().Bool("no-destroy", false, "Do not destroy the session when complete")
	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)

	return cmd
}

func kvUnlock(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	session := viper.GetString("session")
	if session == "" {
		return fmt.Errorf("Session ID must be provided")
	}

	client, err := newKv()
	if err != nil {
		return err
	}

	sessionClient, err := newSession()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	kv, _, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		return fmt.Errorf("Node '%s' does not exist", path)
	}

	if kv.Session != session {
		return fmt.Errorf("Session not lock holder")
	}

	writeOpts := writeOptions()

	success, _, err := client.Release(kv, writeOpts)
	if err != nil {
		return err
	}

	if !viper.GetBool("no-destroy") {
		_, err = sessionClient.Destroy(session, writeOpts)
		if err != nil {
			return err
		}
	}

	if !success {
		return fmt.Errorf("Failed unlocking path")
	}

	return nil
}

// Watch functions

func newKvWatchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch <path>",
		Short: "Watch for changes to a K/V path",
		Long:  "Watch for changes to a K/V path",
		RunE:  kvWatch,
	}

	cmd.Flags().String("fields", "all", "Comma separated list of fields to return.")
	cmd.Flags().String("format", "prettyjson", "Output format. Supported options: text, json, prettyjson")
	cmd.Flags().String("delimited", "", "Output field delimiter")
	cmd.Flags().Bool("header", false, "Output a header row for text format")

	addConsistencyOptions(cmd)
	addDatacenterOption(cmd)
	addWaitIndexOption(cmd)
	addTemplateOption(cmd)

	return cmd
}

func kvWatch(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single key path must be specified")
	}
	path := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newKv()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	kvo := NewKVOutput(os.Stdout, os.Stderr, viper.GetString("fields"))

RETRY:
	kv, meta, err := client.Get(path, queryOpts)
	if err != nil {
		return err
	}

	if kv == nil {
		queryOpts.WaitIndex = meta.LastIndex
		goto RETRY
	}

	if viper.GetString("template") != "" {
		return output(kv)
	} else {
		return kvo.Output(kv, OutputFormat{
			Type:      viper.GetString("format"),
			Delimiter: viper.GetString("delimiter"),
			Header:    viper.GetBool("header"),
		})
	}
}

// Write functions

func newKvWriteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write <path> <value>",
		Short: "Write a value to a given path",
		Long:  "Write a value to a given path",
		RunE:  kvWrite,
	}

	cmd.Flags().String("modifyindex", "", "Perform a Check-and-Set write")
	cmd.Flags().String("flags", "", "Integer value between 0 and 2^64 - 1")
	addDatacenterOption(cmd)

	return cmd
}

func kvWrite(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Key path and value must be specified")
	}

	path := args[0]
	value := strings.Join(args[1:], " ")

	viper.BindPFlags(cmd.Flags())

	kv := new(consulapi.KVPair)

	kv.Key = path
	if strings.HasPrefix(value, "@") {
		v, err := ioutil.ReadFile(value[1:])
		if err != nil {
			return fmt.Errorf("ReadFile error: %v", err)
		}
		kv.Value = v
	} else {
		kv.Value = []byte(value)
	}

	// &flags=
	//
	if flags := viper.GetString("flags"); flags != "" {
		f, err := strconv.ParseUint(flags, 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing flags: %v", flags)
		}
		kv.Flags = f
	}

	client, err := newKv()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	if viper.GetString("modifyindex") == "" {
		_, err := client.Put(kv, writeOpts)
		if err != nil {
			return err
		}
	} else {
		// Check-and-Set
		i, err := strconv.ParseUint(viper.GetString("modifyindex"), 0, 64)
		if err != nil {
			return fmt.Errorf("Error parsing modifyIndex: %v", viper.GetString("modifyindex"))
		}
		kv.ModifyIndex = i

		success, _, err := client.CAS(kv, writeOpts)
		if err != nil {
			return err
		}

		if !success {
			return fmt.Errorf("Failed to write to K/V")
		}
	}

	return nil
}
