package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func Init(name, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "consul-cli",
		Short:         "Command line interface for Consul HTTP API",
		Long:          "Command line interface for Consul HTTP API",
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.HelpFunc()(cmd, []string{})
			return nil
		},
	}

	cmd.PersistentFlags().String("consul", "", "Consul address:port")
	cmd.PersistentFlags().Bool("ssl", false, "Use HTTPS when talking to Consul")
	cmd.PersistentFlags().Bool("ssl-verify", true, "Verify certificates when connecting via SSL")
	cmd.PersistentFlags().String("ssl-cert", "", "Path to an SSL client certificate for authentication")
	cmd.PersistentFlags().String("ssl-key", "", "Path to an SSL client certificate key for authentication")
	cmd.PersistentFlags().String("ssl-ca-cert", "", "Path to a CA certificate file to validate the Consul server")
	cmd.PersistentFlags().String("auth", "", "The HTTP basic authentication username (and optional password) separated by a colon")
	cmd.PersistentFlags().String("token", "", "The Consul ACL token")
	cmd.PersistentFlags().String("token-file", "", "Path to file containing Consul ACL token")
	cmd.PersistentFlags().BoolVarP(&cmd.SilenceUsage, "quiet", "q", true, "Don't show usage on error")

	cmd.AddCommand(newAclCommand())
	cmd.AddCommand(newAgentCommand())
	cmd.AddCommand(newCatalogCommand())
	cmd.AddCommand(newCheckCommand())
	cmd.AddCommand(newCoordinateCommand())
	cmd.AddCommand(newKvCommand())
	cmd.AddCommand(newHealthCommand())
	cmd.AddCommand(newServiceCommand())
	cmd.AddCommand(newSessionCommand())
	cmd.AddCommand(newStatusCommand())

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s %s\n", name, version)
			return nil
		},
	}
	cmd.AddCommand(versionCmd)

	return cmd
}

func addTemplateOption(cmd *cobra.Command) {
	cmd.Flags().String("template", "", "Output template. Use @filename to read template from a file")
}
