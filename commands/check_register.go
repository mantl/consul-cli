package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

	addCheckOptions(cmd)
	addRawOption(cmd)

	cmd.Flags().String("service-id", "", "Service ID to associate check")

	return cmd
}

func checkRegister(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var check consulapi.AgentCheckRegistration

	if raw := viper.GetString("raw"); raw != "" {
		if err := readRawJSON(raw, &check); err != nil {
			return err
		}
	} else {
		if len(args) != 1 {
			return fmt.Errorf("A single check id must be specified")
		}
		checkName := args[0]

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

		check = consulapi.AgentCheckRegistration{
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
	}

	client, err := newAgent()
	if err != nil {
		return err
	}

	err = client.CheckRegister(&check)
	if err != nil {
		return err
	}

	return nil
}
