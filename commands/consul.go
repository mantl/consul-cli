package commands

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
)

type consul struct {
	configFile string
	env        string
	address    string
	sslEnabled bool
	sslVerify  bool
	sslCert    string
	sslCaCert  string
	sslKey     string
	token      string
	tokenFile  string
	auth       *auth
	tlsConfig  *tls.Config

	dc        string
	waitIndex uint64
}

func (c *Cmd) ACL() (*consulapi.ACL, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.ACL(), nil
}

func (c *Cmd) Agent() (*consulapi.Agent, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Agent(), nil
}

func (c *Cmd) Catalog() (*consulapi.Catalog, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Catalog(), nil
}

func (c *Cmd) Coordinate() (*consulapi.Coordinate, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Coordinate(), nil
}

func (c *Cmd) Health() (*consulapi.Health, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Health(), nil
}

func (c *Cmd) KV() (*consulapi.KV, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.KV(), nil
}

func (c *Cmd) Session() (*consulapi.Session, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Session(), nil
}

func (c *Cmd) Status() (*consulapi.Status, error) {
	consul, err := c.Client()
	if err != nil {
		return nil, err
	}

	return consul.Status(), nil
}

// Get the values from ~/.consul-cli which will override
// what is specified in the command-line
func (c *Cmd) Client() (*consulapi.Client, error) {
	ConsulFile := "~/.consul-cli"
	if ConsulFile := CheckConsulFile(ConsulFile); ConsulFile != "" {
		configFromFile := ReadConsulFile(ConsulFile, c.consul.env)
		if _, ok := configFromFile["consul"]; ok {
			c.consul.address = configFromFile["consul"].(string)
		}
		if _, ok := configFromFile["ssl"]; ok {
			c.consul.sslEnabled = configFromFile["ssl"].(bool)
		}
		if _, ok := configFromFile["ssl-verify"]; ok {
			c.consul.sslVerify = configFromFile["ssl-verify"].(bool)
		}
		if _, ok := configFromFile["ssl-cert"]; ok {
			c.consul.sslCert = configFromFile["ssl-cert"].(string)
		}
		if _, ok := configFromFile["ssl-ca-cert"]; ok {
			c.consul.sslCaCert = configFromFile["ssl-ca-cert"].(string)
		}
		if _, ok := configFromFile["ssl-key"]; ok {
			c.consul.sslCert = configFromFile["ssl-key"].(string)
		}
		if _, ok := configFromFile["token"]; ok {
			c.consul.token = configFromFile["token"].(string)
		}
	}

	config := consulapi.DefaultConfig()
	csl := c.consul
	csl.tlsConfig = new(tls.Config)

	if csl.address != "" {
		config.Address = c.consul.address
	}

	if csl.token != "" && csl.tokenFile != "" {
		return nil, errors.New("--token and --token-file can not both be provided")
	}

	if csl.tokenFile != "" {
		b, err := ioutil.ReadFile(csl.tokenFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading token file: %s", err)
		}

		config.Token = strings.TrimSpace(string(b))
	}

	if csl.token != "" {
		config.Token = csl.token
	}

	if csl.sslEnabled {
		config.Scheme = "https"

		if csl.sslCert != "" {
			if csl.sslKey == "" || csl.sslCaCert == "" {
				return nil, errors.New("--ssl-key and --ssl-ca-cert must be provided in order to use certificates for authentication")
			}
			clientCert, err := tls.LoadX509KeyPair(csl.sslCert, csl.sslKey)
			if err != nil {
				return nil, err
			}

			caCert, err := ioutil.ReadFile(csl.sslCaCert)
			if err != nil {
				return nil, err
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			csl.tlsConfig.Certificates = []tls.Certificate{clientCert}
			csl.tlsConfig.RootCAs = caCertPool
			csl.tlsConfig.BuildNameToCertificate()
		}
	}

	transport := new(http.Transport)
	transport.TLSClientConfig = csl.tlsConfig

	if !csl.sslVerify {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
	config.HttpClient.Transport = transport

	if csl.auth.Enabled {
		config.HttpAuth = &consulapi.HttpBasicAuth{
			Username: csl.auth.Username,
			Password: csl.auth.Password,
		}
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Cmd) WriteOptions() *consulapi.WriteOptions {
	csl := c.consul

	writeOpts := new(consulapi.WriteOptions)
	if csl.token != "" {
		writeOpts.Token = csl.token
	}

	if csl.dc != "" {
		writeOpts.Datacenter = csl.dc
	}

	return writeOpts
}

func (c *Cmd) QueryOptions() *consulapi.QueryOptions {
	csl := c.consul

	queryOpts := new(consulapi.QueryOptions)
	if csl.token != "" {
		queryOpts.Token = csl.token
	}

	if csl.dc != "" {
		queryOpts.Datacenter = csl.dc
	}

	if csl.waitIndex != 0 {
		queryOpts.WaitIndex = csl.waitIndex
	}

	return queryOpts
}

func (c *Cmd) AddDatacenterOption(cmd *cobra.Command) {
	cmd.Flags().StringVar(&c.consul.dc, "datacenter", "", "Consul data center")
}

func (c *Cmd) AddWaitIndexOption(cmd *cobra.Command) {
	cmd.Flags().Uint64Var(&c.consul.waitIndex, "wait-index", 0, "Only return if ModifyIndex is greater than <index>")
}

func NewConsul() *consul {
	return &consul{
		auth: new(auth),
	}
}

type auth struct {
	Enabled  bool
	Username string
	Password string
}

func (a *auth) Set(value string) error {
	a.Enabled = true

	if strings.Contains(value, ":") {
		split := strings.SplitN(value, ":", 2)
		a.Username = split[0]
		a.Password = split[1]
	} else {
		a.Username = value
	}

	return nil
}

func (a *auth) String() string {
	if a.Password == "" {
		return a.Username
	}

	return fmt.Sprintf("%s:%s", a.Username, a.Password)
}

func (a *auth) Type() string {
	return "auth"
}
