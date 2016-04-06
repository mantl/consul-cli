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
	address    string
	sslEnabled bool
	sslVerify  bool
	sslCert    string
	sslKey     string
	sslCaCert  string
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

func (c *Cmd) Client() (*consulapi.Client, error) {
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
			if csl.sslKey == "" {
				return nil, errors.New("--ssl-key must be provided in order to use certificates for authentication")
			}
			clientCert, err := tls.LoadX509KeyPair(csl.sslCert, csl.sslKey)
			if err != nil {
				return nil, err
			}

			csl.tlsConfig.Certificates = []tls.Certificate{clientCert}
			csl.tlsConfig.BuildNameToCertificate()
		}

		if csl.sslVerify {
			if csl.sslCaCert == "" {
				return nil, errors.New("--ssl-ca-cert must be provided in order to use certificates for verification")
			}

			caCert, err := ioutil.ReadFile(csl.sslCaCert)
			if err != nil {
				return nil, err
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			csl.tlsConfig.RootCAs = caCertPool
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
