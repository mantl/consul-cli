package commands

import (
	"errors"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newServiceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service",
		Short: "Consul /agent/service endpoint interface",
		Long:  "Consul /agent/service endpoint interface",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, []string{})
			return nil
		},
	}

	cmd.AddCommand(newServiceDeregisterCommand())
	cmd.AddCommand(newServiceMaintenanceCommand())
	cmd.AddCommand(newServiceRegisterCommand())

	return cmd
}

// Deregistration functions

func newServiceDeregisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deregister <serviceId>",
		Short: "Remove a service from the agent",
		Long:  "Remove a service from the agent",
		RunE:  serviceDeregister,
	}

	return cmd
}

func serviceDeregister(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	viper.BindPFlags(cmd.Flags())

	agent, err := newAgent()
	if err != nil {
		return err
	}

	var result error

	for _, id := range args {
		if err := agent.ServiceDeregister(id); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

// Maintenance functions

func newServiceMaintenanceCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "maintenance",
		Short: "Manage maintenance mode of a service",
		Long:  "Manage maintenance mode of a service",
		RunE:  serviceMaintenance,
	}

	cmd.Flags().Bool("enabled", true, "Boolean value for maintenance mode")
	cmd.Flags().String("reason", "", "Reason for entering maintenance mode")

	return cmd
}

func serviceMaintenance(cmd *cobra.Command, args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("No service IDs specified")
	}

	viper.BindPFlags(cmd.Flags())

	agent, err := newAgent()
	if err != nil {
		return err
	}

	var result error

	enabled := viper.GetBool("enabled")
	reason := viper.GetString("reason")

	for _, id := range args {
		if enabled {
			if err := agent.EnableServiceMaintenance(id, reason); err != nil {
				result = multierror.Append(result, err)
			}
		} else {
			if err := agent.DisableServiceMaintenance(id); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	return result
}

// Registration functions

var srLongHelp = `Register a new local service

  If --id is not specified, the serviceName is used. There cannot
be duplicate service IDs per agent however.

  If --address is not specified, the IP address of the local agent
is used.
`

func newServiceRegisterCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <serviceName>",
		Short: "Register a new local service",
		Long:  srLongHelp,
		RunE:  serviceRegister,
	}

	addAgentServiceOptions(cmd)
	addRawOption(cmd)

	return cmd
}

func serviceRegister(cmd *cobra.Command, args []string) error {
	viper.BindPFlags(cmd.Flags())

	var service consulapi.AgentServiceRegistration

	if raw := viper.GetString("raw"); raw != "" {
		if err := readRawJSON(raw, &service); err != nil {
			return err
		}
	} else {
		switch {
		case len(args) == 0:
			return fmt.Errorf("Service name must be specified")
		case len(args) > 1:
			return fmt.Errorf("Only one service name allowed")
		}
		serviceName := args[0]

		//checkStrings := viper.GetStringSlice("check")
		//checks, err := parseChecks(checkStrings)
		//if err != nil {
		//	return err
		//}

		service = consulapi.AgentServiceRegistration{
			ID:                viper.GetString("id"),
			Name:              serviceName,
			Tags:              getStringSlice(cmd, "tag"),
			Port:              viper.GetInt("port"),
			Address:           viper.GetString("address"),
//			Checks:            checks,
			EnableTagOverride: viper.GetBool("override-tag"),
		}
	}

	agent, err := newAgent()
	if err != nil {
		return err
	}

	if err := agent.ServiceRegister(&service); err != nil {
		return err
	}

	return nil
}

func parseChecks(checks []string) ([]*consulapi.AgentServiceCheck, error) {
	rval := make([]*consulapi.AgentServiceCheck, len(checks))

	for i, s := range checks {
		parts := strings.Split(s, ":")
		numParts := len(parts)

		switch parts[0] {
		case "http":
			// Order is http:interval:url:[tlsskip]
			switch {
			case numParts < 3 || numParts > 4:
				return nil, fmt.Errorf("HTTP check format is http:interval:url:[tlsskipverify]")
			case numParts == 3:
				rval[i] = &consulapi.AgentServiceCheck{
					HTTP:     parts[2],
					Interval: parts[1],
				}
			case numParts == 4:
				rval[i] = &consulapi.AgentServiceCheck{
					HTTP:     parts[2],
					Interval: parts[1],
					//TLSSkipVerify: parts[3],
				}
			}
		case "script":
			// Order is script:interval:command
			if numParts != 3 {
				return nil, fmt.Errorf("Script check format is script:interval:command_path")
			}

			rval[i] = &consulapi.AgentServiceCheck{
				Script:   parts[2],
				Interval: parts[1],
			}
		case "ttl":
			// Order is ttl:interval
			if numParts != 2 {
				return nil, fmt.Errorf("TTL check format is ttl:interval")
			}

			rval[i] = &consulapi.AgentServiceCheck{
				TTL: parts[1],
			}
		}
	}

	return rval, nil
}
func parseCheckConfig(s string) (*consulapi.AgentServiceCheck, error) {
	if len(strings.TrimSpace(s)) < 1 {
		return nil, errors.New("Cannot specify empty check")
	}

	var checkType, checkInterval, checkString string
	parts := strings.Split(s, ":")
	partLen := len(parts)
	switch {
	case partLen == 2:
		checkType, checkInterval = parts[0], parts[1]
		checkString = ""
	case partLen >= 3:
		checkType, checkInterval, checkString = parts[0], parts[1], strings.Join(parts[2:], ":")
	default:
		return nil, fmt.Errorf("Invalid check definition '%s'", s)
	}

	switch strings.ToLower(checkType) {
	case "http":
		if checkString == "" {
			return nil, errors.New("Must provide URL with HTTP check")
		}

		return &consulapi.AgentServiceCheck{
			HTTP:     checkString,
			Interval: checkInterval,
		}, nil
	case "script":
		if checkString == "" {
			return nil, errors.New("Must provide command path with script check")
		}

		return &consulapi.AgentServiceCheck{
			Script:   checkString,
			Interval: checkInterval,
		}, nil
	case "ttl":
		return &consulapi.AgentServiceCheck{
			TTL: checkInterval,
		}, nil
	}

	return nil, fmt.Errorf("Invalid check type '%s'", checkType)
}
