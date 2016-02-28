package commands

import (
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Cmd struct {
	root *cobra.Command

	Err io.Writer
	Out io.Writer

	consul *consul

	Template string
}

type ConfigFromFile struct {
	Consul    string
	Ssl       bool   `toml:"ssl"`
	SslCaCert string `toml:"ssl-ca-cert"`
	SslCert   string `toml:"ssl-cert"`
	SslVerify bool   `toml:"ssl-verify"`
	Token     string
}

func CheckConfigFile() string {
	user, _ := user.Current()
	homeDir := user.HomeDir
	fmt.Println(homeDir)
	tempConfigFile := fmt.Sprint(homeDir, "/.consul-cli")
	_, err := os.Stat(tempConfigFile)

	if err == nil {
		return tempConfigFile
	}
	return ""
}

func ReadConfigFile(configFile string) (toml.MetaData, ConfigFromFile) {
	var config ConfigFromFile
	metadata, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		fmt.Println("ERROR")
	}
	return metadata, config
}

func Init(name, version string) *Cmd {
	c := Cmd{
		Err: os.Stderr,
		Out: os.Stdout,
		consul: &consul{
			auth: new(auth),
		},
	}

	c.root = &cobra.Command{
		Use:   "consul-cli",
		Short: "Command line interface for Consul HTTP API",
		Long:  "Command line interface for Consul HTTP API",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c.root.Help()
			return nil
		},
	}

	tempConsul := "127.0.0.1:8500"
	tempSslEnabled := false
	tempSslVerify := true
	tempSslCert := ""
	tempSslCaCert := ""
	tempToken := ""

	if tempConfigFile := CheckConfigFile(); tempConfigFile != "" {
		metadata, configFromFile := ReadConfigFile(tempConfigFile)

		if metadata.IsDefined("consul") {
			tempConsul = configFromFile.Consul
		}
		if metadata.IsDefined("ssl") {
			tempSslEnabled = configFromFile.Ssl
		}
		if metadata.IsDefined("ssl-verify") {
			tempSslVerify = configFromFile.SslVerify
		}
		if metadata.IsDefined("ssl-cert") {
			tempSslCert = configFromFile.SslCert
		}
		if metadata.IsDefined("ssl-ca-cert") {
			tempSslCaCert = configFromFile.SslCaCert
		}
		if metadata.IsDefined("token") {
			tempToken = configFromFile.Token
		}
	}

	c.root.PersistentFlags().StringVar(&c.consul.configFile, "consul-file", "~/.consul-cli", "Configuration file")
	c.root.PersistentFlags().StringVar(&c.consul.address, "consul", tempConsul, "Consul address:port")
	c.root.PersistentFlags().BoolVar(&c.consul.sslEnabled, "ssl", tempSslEnabled, "Use HTTPS when talking to Consul")
	c.root.PersistentFlags().BoolVar(&c.consul.sslVerify, "ssl-verify", tempSslVerify, "Verify certificates when connecting via SSL")
	c.root.PersistentFlags().StringVar(&c.consul.sslCert, "ssl-cert", tempSslCert, "Path to an SSL client certificate for authentication")
	c.root.PersistentFlags().StringVar(&c.consul.sslCaCert, "ssl-ca-cert", tempSslCaCert, "Path to a CA certificate file to validate the Consul server")
	c.root.PersistentFlags().Var((*auth)(c.consul.auth), "auth", "The HTTP basic authentication username (and optional password) separated by a colon")
	c.root.PersistentFlags().StringVar(&c.consul.token, "token", tempToken, "The Consul ACL token")

	c.initAcl()
	c.initAgent()
	c.initCatalog()
	c.initCheck()
	c.initCoordinate()
	c.initHealth()
	c.initKv()
	c.initService()
	c.initSession()
	c.initStatus()

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("%s %s\n", name, version)
			return nil
		},
	}
	c.root.AddCommand(versionCmd)

	return &c
}

func (c *Cmd) Execute() error {
	return c.root.Execute()
}

func (c *Cmd) Help() error {
	return c.root.Help()
}

func (c *Cmd) AddCommand(cmd *cobra.Command) {
	c.root.AddCommand(cmd)
}

func (c *Cmd) AddTemplateOption(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.Template, "template", "", "Output template. Use @filename to read template from a file")
}

type funcVar func(s string) error

func (f funcVar) Set(s string) error { return f(s) }
func (f funcVar) String() string     { return "" }
func (f funcVar) IsBoolFlag() bool   { return false }
func (f funcVar) Type() string       { return "funcVar" }
