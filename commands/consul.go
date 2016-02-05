package commands

import (
	"crypto/tls"
	"fmt"
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
	sslCaCert  string
	token      string
	auth       *auth

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

	if csl.address != "" {
		config.Address = c.consul.address
	}

	if csl.token != "" {
		config.Token = csl.token
	}

	if csl.sslEnabled {
		config.Scheme = "https"
	}

	if !csl.sslVerify {
		config.HttpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

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
