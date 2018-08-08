package action

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	consulapi "github.com/hashicorp/consul/api"
)

type consul struct {
	address   string
	dc        string
	ssl       bool
	sslVerify bool
	sslCert   string
	sslKey    string
	sslCaCert string
	auth      string
	token     string
	tokenFile string

	waitIndex  uint64
	consistent bool
	stale      bool
	nodeMeta   []string
	near       string
}

func (c *consul) newACL() (*consulapi.ACL, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.ACL(), nil
}

func (c *consul) newAgent() (*consulapi.Agent, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Agent(), nil
}

func (c *consul) newCatalog() (*consulapi.Catalog, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Catalog(), nil
}

func (c *consul) newCoordinate() (*consulapi.Coordinate, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Coordinate(), nil
}

func (c *consul) newEvent() (*consulapi.Event, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Event(), nil
}

func (c *consul) newHealth() (*consulapi.Health, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Health(), nil
}

func (c *consul) newKv() (*consulapi.KV, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.KV(), nil
}

func (c *consul) newOperator() (*consulapi.Operator, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Operator(), nil
}

func (c *consul) newSession() (*consulapi.Session, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Session(), nil
}

func (c *consul) newSnapshot() (*consulapi.Snapshot, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Snapshot(), nil
}

func (c *consul) newStatus() (*consulapi.Status, error) {
	consul, err := c.newClient()
	if err != nil {
		return nil, err
	}

	return consul.Status(), nil
}

func (c *consul) newClient() (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()

	// Consul address handling

	if c.address != "" {
		config.Address = c.address
	}

	// ACL token handling
	if c.token != "" && c.tokenFile != "" {
		return nil, errors.New("--token and --token-file can not both be provided")
	}

	if c.tokenFile != "" {
		b, err := ioutil.ReadFile(c.tokenFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading token file: %s", err)
		}
		config.Token = strings.TrimSpace(string(b))
	}

	if c.token != "" {
		config.Token = c.token
	}

	// SSL handling

	if c.ssl {
		config.Scheme = "https"

		tlsConfig := new(tls.Config)
		transport := new(http.Transport)

		if c.sslCert != "" {
			if c.sslKey == "" {
				return nil, errors.New("--ssl-key must be provided in order to use certificates for authentication")
			}
			clientCert, err := tls.LoadX509KeyPair(c.sslCert, c.sslKey)
			if err != nil {
				return nil, err
			}

			tlsConfig.Certificates = []tls.Certificate{clientCert}
			tlsConfig.BuildNameToCertificate()
		}

		if c.sslVerify {
			if c.sslCaCert == "" {
				return nil, errors.New("--ssl-ca-cert must be provided in order to use certificates for verification")
			}

			caCert, err := ioutil.ReadFile(c.sslCaCert)
			if err != nil {
				return nil, err
			}

			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		} else {
			tlsConfig.InsecureSkipVerify = true
		}

		transport.TLSClientConfig = tlsConfig
		config.Transport = transport
	}

	// Auth handling

	if c.auth != "" {
		auth := new(consulapi.HttpBasicAuth)
		if strings.Contains(c.auth, ":") {
			split := strings.SplitN(c.auth, ":", 2)
			auth.Username = split[0]
			auth.Password = split[1]
		} else {
			auth.Username = c.auth
		}
		config.HttpAuth = auth
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, err
}

func (c *consul) writeOptions() *consulapi.WriteOptions {
	writeOpts := new(consulapi.WriteOptions)

	if c.token != "" {
		writeOpts.Token = c.token
	}

	if c.dc != "" {
		writeOpts.Datacenter = c.dc
	}

	return writeOpts
}

func (c *consul) queryOptions() *consulapi.QueryOptions {
	queryOpts := new(consulapi.QueryOptions)

	if c.token != "" {
		queryOpts.Token = c.token
	}

	if c.dc != "" {
		queryOpts.Datacenter = c.dc
	}

	if c.waitIndex != 0 {
		queryOpts.WaitIndex = c.waitIndex
	}

	if len(c.nodeMeta) > 0 {
		for _, kvp := range c.nodeMeta {
			parts := strings.Split(kvp, ":")
			if len(parts) == 2 {
				queryOpts.NodeMeta[parts[0]] = parts[1]
			}
		}
	}

	if c.near != "" {
		queryOpts.Near = c.near
	}

	queryOpts.RequireConsistent = c.consistent
	queryOpts.AllowStale = c.stale

	return queryOpts
}

