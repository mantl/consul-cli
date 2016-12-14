package commands

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newOperatorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "operator",
		Short: "Consul /operator endpoint interface",
		Long:  "Consul /operator endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorKeyringCommand())
	cmd.AddCommand(newOperatorRaftCommand())

	return cmd
}

// keyring command

func newOperatorKeyringCommand() *cobra.Command {
	cmd := &cobra.Command{
		Hidden: true, // Hide subcommand Consul official release
		Use:    "keyring",
		Short:  "Consul /operator/keyring interface",
		Long:   "Consul /operator/keyring interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorKeyringInstallCommand())
	cmd.AddCommand(newOperatorKeyringListCommand())
	cmd.AddCommand(newOperatorKeyringRemoveCommand())
	cmd.AddCommand(newOperatorKeyringUseCommand())

	return cmd
}

// keyring install functions

func newOperatorKeyringInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <key> [<key>]",
		Short: "Install a new gossip key into the cluster",
		Long:  "Install a new gossip key into the cluster",
		RunE:  operatorKeyringInstall,
	}

	return cmd
}

func operatorKeyringInstall(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	if len(args) < 1 {
		return fmt.Errorf("At least one gossip key must be specified")
	}

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	var result error

	for _, k := range args {
		if err := client.KeyringInstall(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

// keyring list functions

func newOperatorKeyringListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List gossip keys installed",
		Long:  "List gossip keys installed",
		RunE:  operatorKeyringList,
	}

	return cmd
}

func operatorKeyringList(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()
	r, err := client.KeyringList(queryOpts)
	if err != nil {
		return err
	}

	return output(r)
}

// keyring removel functions

func newOperatorKeyringRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <key> [<key>]",
		Short: "Remove gossip keys from the cluster",
		Long:  "Remove gossip keys from the cluster",
		RunE:  operatorKeyringRemove,
	}

	return cmd
}

func operatorKeyringRemove(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	if len(args) < 1 {
		return fmt.Errorf("At least one gossip key must be specified")
	}

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	var result error

	for _, k := range args {
		if err := client.KeyringRemove(k, writeOpts); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

// keyring use functions

func newOperatorKeyringUseCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use <key>",
		Short: "Change the primary gossip encryption key",
		Long:  "Change the primary gossip encryption key",
		RunE:  operatorKeyringUse,
	}

	return cmd
}

func operatorKeyringUse(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Encryption key must be specified")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	return client.KeyringUse(args[0], writeOpts)
}

// raft command

func newOperatorRaftCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "raft",
		Short: "Consul /operator/raft endpoint interface",
		Long:  "Consul /operator/raft endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newOperatorRaftConfigCommand())
	cmd.AddCommand(newOperatorRaftDeleteCommand())

	return cmd
}

// Raft configuration functions

func newOperatorRaftConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Inspect the Raft configuration",
		Long:  "Inspect the Raft configuration",
		RunE:  operatorRaftConfig,
	}

	addTemplateOption(cmd)
	addDatacenterOption(cmd)

	cmd.Flags().Bool("stale", false, "Read the raft configuration from any Consul server")

	return cmd
}

func operatorRaftConfig(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	rc, err := client.RaftGetConfiguration(queryOpts)
	if err != nil {
		return err
	}

	return output(rc)
}

// Peer delete functions

var raftDeleteLongHelp = `
Remove a Consul server from the Raft configuration

An address is required and should be set to IP:port for the server
to remove. The port number is 8300 unless configured otherwise`

func newOperatorRaftDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <address>",
		Short: "Remove a Consul server from the Raft configuration",
		Long:  raftDeleteLongHelp,
		RunE:  operatorRaftDelete,
	}

	addDatacenterOption(cmd)

	return cmd
}

func operatorRaftDelete(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single address argument must be specified")
	}
	address := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newOperator()
	if err != nil {
		return err
	}

	writeOpts := writeOptions()

	return client.RaftRemovePeerByAddress(address, writeOpts)
}
