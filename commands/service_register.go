package commands

import (
	"errors"
	"fmt"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type ServiceRegisterOptions struct {
	Id      string
	Tags    []string
	Address string
	Port    int
	Checks  consulapi.AgentServiceChecks
}

var srLongHelp = `Register a new local service

  If --id is not specified, the serviceName is used. There cannot
be duplicate service IDs per agent however.

  If --address is not specified, the IP address of the local agent
is used.
`

func (s *Service) AddRegisterSub(cmd *cobra.Command) {
	sro := &ServiceRegisterOptions{}

	registerCmd := &cobra.Command{
		Use:   "register <serviceName>",
		Short: "Register a new local service",
		Long:  srLongHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Register(args, sro)
		},
	}

	oldRegisterCmd := &cobra.Command{
		Use:        "service-register <serviceName>",
		Short:      "Register a new local service",
		Long:       srLongHelp,
		Deprecated: "Use service register",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return s.Register(args, sro)
		},
	}

	registerCmd.Flags().StringVar(&sro.Id, "id", "", "Service Id")
	registerCmd.Flags().Var((funcVar)(func(s string) error {
		if sro.Tags == nil {
			sro.Tags = make([]string, 0, 1)
		}
		sro.Tags = append(sro.Tags, s)
		return nil
	}), "tag", "Service tag. Multiple tags allowed")
	registerCmd.Flags().StringVar(&sro.Address, "address", "", "Service address")
	registerCmd.Flags().IntVar(&sro.Port, "port", 0, "Service port")
	registerCmd.Flags().Var((funcVar)(func(s string) error {
		t, err := ParseCheckConfig(s)
		if err != nil {
			return err
		}

		if sro.Checks == nil {
			sro.Checks = make(consulapi.AgentServiceChecks, 0, 1)
		}

		sro.Checks = append(sro.Checks, t)
		return nil
	}), "check", "Check to create. Multiple checks can be created with the service. The format is: [http|script|ttl]:<interval>:<command|url>")

	oldRegisterCmd.Flags().StringVar(&sro.Id, "id", "", "Service Id")
	oldRegisterCmd.Flags().Var((funcVar)(func(s string) error {
		if sro.Tags == nil {
			sro.Tags = make([]string, 0, 1)
		}
		sro.Tags = append(sro.Tags, s)
		return nil
	}), "tag", "Service tag. Multiple tags allowed")
	oldRegisterCmd.Flags().StringVar(&sro.Address, "address", "", "Service address")
	oldRegisterCmd.Flags().IntVar(&sro.Port, "port", 0, "Service port")
	oldRegisterCmd.Flags().Var((funcVar)(func(s string) error {
		t, err := ParseCheckConfig(s)
		if err != nil {
			return err
		}

		if sro.Checks == nil {
			sro.Checks = make(consulapi.AgentServiceChecks, 0, 1)
		}

		sro.Checks = append(sro.Checks, t)
		return nil
	}), "check", "Check to create. Multiple checks can be created with the service. The format is: [http|script|ttl]:<interval>:<command|url>")

	cmd.AddCommand(registerCmd)

	s.AddCommand(oldRegisterCmd)
}

func (s *Service) Register(args []string, sro *ServiceRegisterOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Service name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only on service name allowed")
	}
	serviceName := args[0]

	consul, err := s.Client()
	if err != nil {
		return err
	}

	service := &consulapi.AgentServiceRegistration{
		ID:      sro.Id,
		Name:    serviceName,
		Tags:    sro.Tags,
		Port:    sro.Port,
		Address: sro.Address,
		Checks:  sro.Checks,
	}

	client := consul.Agent()
	err = client.ServiceRegister(service)
	if err != nil {
		return err
	}

	return nil
}

func ParseCheckConfig(s string) (*consulapi.AgentServiceCheck, error) {
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
