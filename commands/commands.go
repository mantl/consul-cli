package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

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

type TomlConfig struct {
	Env   string
	Dev   ConsulFromFile
	Prod  ConsulFromFile
	Stage ConsulFromFile
	QA    ConsulFromFile
	West  ConsulFromFile
	East  ConsulFromFile
}

type ConsulFromFile struct {
	Consul    string
	Ssl       bool   `toml:"ssl"`
	SslCaCert string `toml:"ssl-ca-cert"`
	SslCert   string `toml:"ssl-cert"`
	SslKey    string `toml:"ssl-key"`
	SslVerify bool   `toml:"ssl-verify"`
	Token     string
}

func CheckConsulFile(inputFile string) string {
	tempConsulFile := inputFile
	if inputFile[:2] == "~/" {
		user, _ := user.Current()
		homeDir := fmt.Sprint(user.HomeDir + "/")
		tempConsulFile = strings.Replace(inputFile, "~/", homeDir, 1)
	}
	_, err := os.Stat(tempConsulFile)
	if err == nil {
		return tempConsulFile
	}
	return ""
}

func ReadConsulFile(configFile, env string) map[string]interface{} {
	var config TomlConfig
	_, err := toml.DecodeFile(configFile, &config)
	if err != nil {
		log.Fatal(err)
	}

	realConfig := map[string]interface{}{}
	if env == "dev" {
		realConfig["consul"] = config.Dev.Consul
		realConfig["ssl"] = config.Dev.Ssl
		realConfig["ssl-ca-cert"] = config.Dev.SslCaCert
		realConfig["ssl-cert"] = config.Dev.SslCert
		realConfig["ssl-key"] = config.Dev.SslKey
		realConfig["ssl-verify"] = config.Dev.SslVerify
		realConfig["token"] = config.Dev.Token
	}
	if env == "qa" {
		realConfig["consul"] = config.QA.Consul
		realConfig["ssl"] = config.QA.Ssl
		realConfig["ssl-ca-cert"] = config.QA.SslCaCert
		realConfig["ssl-cert"] = config.QA.SslCert
		realConfig["ssl-key"] = config.QA.SslKey
		realConfig["ssl-verify"] = config.QA.SslVerify
		realConfig["token"] = config.QA.Token
	}
	if env == "stage" {
		realConfig["consul"] = config.Stage.Consul
		realConfig["ssl"] = config.Stage.Ssl
		realConfig["ssl-ca-cert"] = config.Stage.SslCaCert
		realConfig["ssl-cert"] = config.Stage.SslCert
		realConfig["ssl-key"] = config.Stage.SslKey
		realConfig["ssl-verify"] = config.Stage.SslVerify
		realConfig["token"] = config.Stage.Token
	}
	if env == "prod" {
		realConfig["consul"] = config.Prod.Consul
		realConfig["ssl"] = config.Prod.Ssl
		realConfig["ssl-ca-cert"] = config.Prod.SslCaCert
		realConfig["ssl-cert"] = config.Prod.SslCert
		realConfig["ssl-key"] = config.Prod.SslKey
		realConfig["ssl-verify"] = config.Prod.SslVerify
		realConfig["token"] = config.Prod.Token
	}
	if env == "west" {
		realConfig["consul"] = config.West.Consul
		realConfig["ssl"] = config.West.Ssl
		realConfig["ssl-ca-cert"] = config.West.SslCaCert
		realConfig["ssl-cert"] = config.West.SslCert
		realConfig["ssl-key"] = config.West.SslKey
		realConfig["ssl-verify"] = config.West.SslVerify
		realConfig["token"] = config.West.Token
	}
	if env == "east" {
		realConfig["consul"] = config.East.Consul
		realConfig["ssl"] = config.East.Ssl
		realConfig["ssl-ca-cert"] = config.East.SslCaCert
		realConfig["ssl-cert"] = config.East.SslCert
		realConfig["ssl-key"] = config.East.SslKey
		realConfig["ssl-verify"] = config.East.SslVerify
		realConfig["token"] = config.East.Token
	}

	return realConfig
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

	c.root.PersistentFlags().StringVar(&c.consul.env, "env", "qa", "default environment")
	c.root.PersistentFlags().StringVar(&c.consul.address, "consul", "127.0.0.1:8500", "Consul address:port")
	c.root.PersistentFlags().BoolVar(&c.consul.sslEnabled, "ssl", false, "Use HTTPS when talking to Consul")
	c.root.PersistentFlags().BoolVar(&c.consul.sslVerify, "ssl-verify", true, "Verify certificates when connecting via SSL")
	c.root.PersistentFlags().StringVar(&c.consul.sslCert, "ssl-cert", "", "Path to an SSL client certificate for authentication")
	c.root.PersistentFlags().StringVar(&c.consul.sslKey, "ssl-key", "", "Path to an SSL client certificate key for authentication")
	c.root.PersistentFlags().StringVar(&c.consul.sslCaCert, "ssl-ca-cert", "", "Path to a CA certificate file to validate the Consul server")
	c.root.PersistentFlags().Var((*auth)(c.consul.auth), "auth", "The HTTP basic authentication username (and optional password) separated by a colon")
	c.root.PersistentFlags().StringVar(&c.consul.token, "token", "", "The Consul ACL token")
	c.root.PersistentFlags().StringVar(&c.consul.tokenFile, "token-file", "", "Path to file containing Consul ACL token")

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
