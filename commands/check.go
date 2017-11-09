package commands

import (
	"github.com/spf13/cobra"

	"github.com/mantl/consul-cli/action"
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
	cmd.AddCommand(newCheckUpdateCommand())
	cmd.AddCommand(newCheckWarnCommand())

	return cmd
}

func newCheckDeregisterCommand() *cobra.Command {
	c := action.CheckDeregisterAction()

	cmd := &cobra.Command{
		Use:   "deregister",
		Short: "Remove a check from the agent",
		Long:  "Remove a check from the agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCheckFailCommand() *cobra.Command {
	c := action.CheckFailAction()

	cmd := &cobra.Command{
		Use:   "fail <checkId>",
		Short: "Mark a local check as critical",
		Long:  "Mark a local check as critical",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}
func newCheckPassCommand() *cobra.Command {
	c := action.CheckPassAction()

	cmd := &cobra.Command{
		Use:   "pass <checkId>",
		Short: "Mark a local check as passing",
		Long:  "Mark a local check as passing",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

var registerLongHelp = `Register a new local check

  If --id is not specified, the checkName is used. There cannot\
be duplicate IDs per agent however.

  Only one of --http, --script, --tcp and --ttl can be specified.
`

func newCheckRegisterCommand() *cobra.Command {
	c := action.CheckRegisterAction()

	cmd := &cobra.Command{
		Use:   "register <checkName>",
		Short: "Register a new local check",
		Long:  registerLongHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCheckUpdateCommand() *cobra.Command {
	c := action.CheckUpdateAction()

	cmd := &cobra.Command{
		Use:   "update <checkId>",
		Short: "Set the status and output of a TTL check",
		Long:  "Set the status and output of a TTL check",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}

func newCheckWarnCommand() *cobra.Command {
	c := action.CheckWarnAction()

	cmd := &cobra.Command{
		Use:   "warn <checkId>",
		Short: "Mark a local check as warning",
		Long:  "Mark a local check as warning",
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(args)
		},
	}

	cmd.Flags().AddGoFlagSet(c.CommandFlags())

	return cmd
}
