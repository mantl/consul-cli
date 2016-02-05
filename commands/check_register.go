package commands

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type CheckRegisterOptions struct {
	Id        string
	Script    string
	Http      string
	Ttl       string
	Interval  string
	Notes     string
	ServiceId string
}

var longHelp = `Register a new local check

  If --id is not specified, the checkName is used. There cannot\
be duplicate IDs per agent however.

  Only one of --http, --script and --ttl can be specified.
`

func (c *Check) AddRegisterSub(cmd *cobra.Command) {
	cro := &CheckRegisterOptions{}

	registerCmd := &cobra.Command{
		Use:   "register <checkName>",
		Short: "Register a new local check",
		Long:  longHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Register(args, cro)
		},
	}
	oldRegisterCmd := &cobra.Command{
		Use:        "check-register <checkName>",
		Short:      "Register a new local check",
		Long:       longHelp,
		Deprecated: "Use acl register",
		Hidden:     true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Register(args, cro)
		},
	}

	registerCmd.Flags().StringVar(&cro.Id, "id", "", "Service Id")
	registerCmd.Flags().StringVar(&cro.Http, "http", "", "A URL to GET every interval")
	registerCmd.Flags().StringVar(&cro.Script, "script", "", "A script to run every interval")
	registerCmd.Flags().StringVar(&cro.Ttl, "ttl", "", "Fail if TTL expires before service checks in")
	registerCmd.Flags().StringVar(&cro.Interval, "interval", "", "Interval between checks")
	registerCmd.Flags().StringVar(&cro.ServiceId, "service-id", "", "Service ID to associate check")
	registerCmd.Flags().StringVar(&cro.Notes, "notes", "", "Description of the check")

	oldRegisterCmd.Flags().StringVar(&cro.Id, "id", "", "Service Id")
	oldRegisterCmd.Flags().StringVar(&cro.Http, "http", "", "A URL to GET every interval")
	oldRegisterCmd.Flags().StringVar(&cro.Script, "script", "", "A script to run every interval")
	oldRegisterCmd.Flags().StringVar(&cro.Ttl, "ttl", "", "Fail if TTL expires before service checks in")
	oldRegisterCmd.Flags().StringVar(&cro.Interval, "interval", "", "Interval between checks")
	oldRegisterCmd.Flags().StringVar(&cro.ServiceId, "service-id", "", "Service ID to associate check")
	oldRegisterCmd.Flags().StringVar(&cro.Notes, "notes", "", "Description of the check")

	cmd.AddCommand(registerCmd)

	c.AddCommand(oldRegisterCmd)
}

func (c *Check) Register(args []string, cro *CheckRegisterOptions) error {
	switch {
	case len(args) == 0:
		return fmt.Errorf("Check name must be specified")
	case len(args) > 1:
		return fmt.Errorf("Only one check name allowed")
	}
	checkName := args[0]

	checkCount := 0
	if cro.Http != "" {
		checkCount = checkCount + 1
	}
	if cro.Script != "" {
		checkCount = checkCount + 1
	}
	if cro.Ttl != "" {
		checkCount = checkCount + 1
	}

	if checkCount > 1 {
		return fmt.Errorf("Only one of --http, --script or --ttl can be specified")
	}

	client, err := c.Agent()
	if err != nil {
		return err
	}

	check := &consulapi.AgentCheckRegistration{
		ID:        cro.Id,
		Name:      checkName,
		ServiceID: cro.ServiceId,
		Notes:     cro.Notes,
		AgentServiceCheck: consulapi.AgentServiceCheck{
			Script:   cro.Script,
			HTTP:     cro.Http,
			Interval: cro.Interval,
			TTL:      cro.Ttl,
		},
	}

	err = client.CheckRegister(check)
	if err != nil {
		return err
	}

	return nil
}
