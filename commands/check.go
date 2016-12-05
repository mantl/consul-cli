package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Consul /agent/check interface",
		Long:  "Consul /agent/check interface",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, []string{})
		},
	}

	cmd.AddCommand(newCheckDeregisterCommand())
	cmd.AddCommand(newCheckFailCommand())
	cmd.AddCommand(newCheckPassCommand())
	cmd.AddCommand(newCheckRegisterCommand())
	cmd.AddCommand(newCheckWarnCommand())

	return cmd
}

// Deregister functions

func newCheckDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister",
		Short: "Remove a check from the agent",
		Long:  "Remove a check from the agent",
		RunE:  checkDeregister,
	}

	return cmd
}

func checkDeregister(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.CheckDeregister(checkId)
	if err != nil {
		return err
	}

	return nil
}

// Fail functions

func newCheckFailCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fail <checkId>",
		Short: "Mark a local check as critical",
		Long:  "Mark a local check as critical",
		RunE:  checkFail,
	}

	cmd.Flags().String("note", "", "Message to associate with check status")

	return cmd
}

func checkFail(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.FailTTL(checkId, viper.GetString("note"))
	if err != nil {
		return err
	}

	return nil
}

// Pass functions

func newCheckPassCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pass <checkId>",
		Short: "Mark a local check as passing",
		Long:  "Mark a local check as passing",
		RunE:  checkPass,
	}

	cmd.Flags().String("note", "", "Message to associate with check status")

	return cmd
}

func checkPass(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.PassTTL(checkId, viper.GetString("note"))
	if err != nil {
		return err
	}

	return nil
}

// Register functions

var registerLongHelp = `Register a new local check

  If --id is not specified, the checkName is used. There cannot\
be duplicate IDs per agent however.

  Only one of --http, --script and --ttl can be specified.
`

func newCheckRegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <checkName>",
		Short: "Register a new local check",
		Long:  registerLongHelp,
		RunE:  checkRegister,
	}

	cmd.Flags().String("id", "", "Service Id")
	cmd.Flags().String("http", "", "A URL to GET every interval")
	cmd.Flags().String("script", "", "A script to run every interval")
	cmd.Flags().String("ttl", "", "Fail if TTL expires before service checks in")
	cmd.Flags().String("interval", "", "Interval between checks")
	cmd.Flags().String("service-id", "", "Service ID to associate check")
	cmd.Flags().String("notes", "", "Description of the check")
	cmd.Flags().String("docker-id", "", "Docker container ID")
	cmd.Flags().String("shell", "", "Shell to use inside docker container")
	cmd.Flags().String("deregister-crit", "", "Deregister critical service after this interval")
	cmd.Flags().Bool("skip-verify", false, "Skip TLS verification for HTTP checks")

	return cmd
}

func checkRegister(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkName := args[0]

	viper.BindPFlags(cmd.Flags())

	checkCount := 0
	if viper.GetString("http") != "" {
		checkCount = checkCount + 1
	}
	if viper.GetString("script") != "" {
		checkCount = checkCount + 1
	}
	if viper.GetString("ttl") != "" {
		checkCount = checkCount + 1
	}

	if checkCount > 1 {
		return fmt.Errorf("Only one of --http, --script or --ttl can be specified")
	}

	client, err := newAgent()
	if err != nil {
		return err
	}

	check := &consulapi.AgentCheckRegistration{
		ID:        viper.GetString("id"),
		Name:      checkName,
		ServiceID: viper.GetString("service-id"),
		Notes:     viper.GetString("notes"),
		AgentServiceCheck: consulapi.AgentServiceCheck{
			Script:            viper.GetString("script"),
			HTTP:              viper.GetString("http"),
			Interval:          viper.GetString("interval"),
			TTL:               viper.GetString("ttl"),
			TLSSkipVerify:     viper.GetBool("skip-verify"),
			DockerContainerID: viper.GetString("docker-id"),
			Shell:             viper.GetString("shell"),
			Notes:             viper.GetString("notes"),
			DeregisterCriticalServiceAfter: viper.GetString("deregister-crit"),
		},
	}

	err = client.CheckRegister(check)
	if err != nil {
		return err
	}

	return nil
}

// Warn functions

func newCheckWarnCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "warn <checkId>",
		Short: "Mark a local check as warning",
		Long:  "Mark a local check as warning",
		RunE:  checkWarn,
	}

	cmd.Flags().String("note", "", "Message to associate with check status")

	return cmd
}

func checkWarn(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("A single check id must be specified")
	}
	checkId := args[0]

	viper.BindPFlags(cmd.Flags())

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.WarnTTL(checkId, viper.GetString("note"))
	if err != nil {
		return err
	}

	return nil
}
