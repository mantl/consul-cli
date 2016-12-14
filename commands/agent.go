package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newAgentCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Consul /agent endpoint interface",
		Long:  "Consul /agent endpoint interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newAgentChecksCommand())
	cmd.AddCommand(newAgentForceLeaveCommand())
	cmd.AddCommand(newAgentJoinCommand())
	cmd.AddCommand(newAgentLeaveCommand())
	cmd.AddCommand(newAgentMaintenanceCommand())
	cmd.AddCommand(newAgentMembersCommand())
	cmd.AddCommand(newAgentMonitorCommand())
	cmd.AddCommand(newAgentReloadCommand())
	cmd.AddCommand(newAgentSelfCommand())
	cmd.AddCommand(newAgentServicesCommand())

	return cmd
}

// Checks functions

func newAgentChecksCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checks",
		Short: "Get the checks the agent is managing",
		Long:  "Get the checks the agent is managing",
		RunE:  agentChecks,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentChecks(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Checks()
	if err != nil {
		return err
	}

	return output(config)
}

// Force Leave functions

func newAgentForceLeaveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "force-leave <node name>",
		Short: "Force the removal of a node",
		Long:  "Force the removal of a node",
		RunE:  agentForceLeave,
	}

	return cmd
}

func agentForceLeave(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Name not provided")
	case len(args) > 1:
		return fmt.Errorf("Only one node allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.ForceLeave(args[0])
}

// Join functions

func newAgentJoinCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "join",
		Short: "Trigger the local agent to join a node",
		Long:  "Trigger the local agent to join a node",
		RunE:  agentJoin,
	}

	cmd.Flags().Bool("wan", false, "Get list of WAN join instead of LAN")

	return cmd
}

func agentJoin(cmd *cobra.Command, args []string) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Node name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one name allowed")
	}

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Join(args[0], viper.GetBool("wan"))
}

// Leave functions

func newAgentLeaveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "leave",
		Short: "Cause the agent to gracefully shutdown and leave the cluster",
		Long:  "Cause the agent to gracefully shutdown and leave the cluster",
		RunE:  agentLeave,
	}

	return cmd
}

func agentLeave(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Leave()
}

// Maintenance functions

func newAgentMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage node maintenance mode",
		Long:  "Manage node maintenance mode",
		RunE:  agentMaintenance,
	}

	cmd.Flags().Bool("enabled", true, "Boolean value for maintenance mode")
	cmd.Flags().String("reason", "", "Reason for entering maintenance mode")

	return cmd
}

func agentMaintenance(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	if viper.GetBool("enabled") {
		return client.EnableNodeMaintenance(viper.GetString("reason"))
	} else {
		return client.DisableNodeMaintenance()
	}
}

// Members functions

func newAgentMembersCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "members",
		Short: "Get the members as seen by the serf agent",
		Long:  "Get the members as seen by the serf agent",
		RunE:  agentMembers,
	}

	cmd.Flags().Bool("wan", false, "Get list of WAN members instead of LAN")

	addTemplateOption(cmd)

	return cmd
}

func agentMembers(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	ms, err := client.Members(viper.GetBool("wan"))
	if err != nil {
		return err
	}

	return output(ms)
}

// Monitor functions

func newAgentMonitorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitor",
		Short: "Streams logs from the agent",
		Long:  "Streams logs from the agent",
		RunE:  agentMonitor,
	}

	cmd.Flags().String("loglevel", "", "Log level to filter on. Default is info")

	return cmd
}

func agentMonitor(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	outputChan, err := client.Monitor(
		viper.GetString("loglevel"),
		nil, // XXX - Set up interrupts and stop things gracefully
		nil, // query options. No documentation on what they are used for. Leave nil for now.
	)
	if err != nil {
		return err
	}

	for s := range outputChan {
		fmt.Printf(s)
	}

	return nil
}

// Reload functions

func newAgentReloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Tell the Consul agent to reload its configuration",
		Long:  "Tell the Consul agent to reload its configuration",
		RunE:  agentReload,
	}

	return cmd
}

func agentReload(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	return client.Reload()
}

// Self functions

func newAgentSelfCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "self",
		Short: "Get agent configuration",
		Long:  "Get agent configuration",
		RunE:  agentSelf,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentSelf(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Self()
	if err != nil {
		return err
	}

	return output(config)
}

// Services functions

func newAgentServicesCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Get the services the agent is managing",
		Long:  "Get the services the agent is managing",
		RunE:  agentServices,
	}

	addTemplateOption(cmd)

	return cmd
}

func agentServices(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	config, err := client.Services()
	if err != nil {
		return err
	}

	return output(config)
}
