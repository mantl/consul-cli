package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newHealthCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Consul /health endpoint interface",
		Long:  "Consul /health endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newHealthChecksCommand())
	cmd.AddCommand(newHealthNodeCommand())
	cmd.AddCommand(newHealthServiceCommand())
	cmd.AddCommand(newHealthStateCommand())

	return cmd
}

// Checks functions

func newHealthChecksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks <serviceName>",
		Short: "Get the health checks for a service",
		Long:  "Get the health checks for a service",
		RunE:  healthChecks,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthChecks(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	checks, _, err := client.Checks(service, queryOpts)
	if err != nil {
		return err
	}

	return output(checks)
}

// Node functions

func newHealthNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "node <nodeName>",
		Short: "Get the health info for a node",
		Long:  "Get the health info for a node",
		RunE:  healthNode,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthNode(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one node name allowed")
	}
	node := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	n, _, err := client.Node(node, queryOpts)
	if err != nil {
		return err
	}

	return output(n)
}

// Service functions

func newHealthServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service <serviceName>",
		Short: "Get the nodes and health info for a service",
		Long:  "Get the nodes and health info for a service",
		RunE:  healthService,
	}

	cmd.Flags().String("tag", "", "Service tag to filter on")
	cmd.Flags().Bool("passing", false, "Only return passing checks")
	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthService(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one service name allowed")
	}
	service := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	tag := viper.GetString("tag")
	po := viper.GetBool("passing")

	s, _, err := client.Service(service, tag, po, queryOpts)
	if err != nil {
		return err
	}

	return output(s)
}

// State functions

func newHealthStateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Get the checks in a given state",
		Long:  "Get the checks in a given state",
		RunE:  healthState,
	}

	addDatacenterOption(cmd)
	addTemplateOption(cmd)
	addConsistencyOptions(cmd)

	return cmd
}

func healthState(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Check state must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one check state allowed")
	}
	state := strings.ToLower(args[0])

	viper.BindPFlags(cmd.Flags())

	client, err := newHealth()
	if err != nil {
		return err
	}

	queryOpts := queryOptions()

	s, _, err := client.State(state, queryOpts)
	if err != nil {
		return err
	}

	return output(s)
}
