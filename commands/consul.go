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
	"github.com/spf13/viper"
)

func newACL() (*consulapi.ACL, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.ACL(), nil
}

func newAgent() (*consulapi.Agent, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Agent(), nil
}

func newCatalog() (*consulapi.Catalog, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Catalog(), nil
}

func newCoordinate() (*consulapi.Coordinate, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Coordinate(), nil
}

func newHealth() (*consulapi.Health, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Health(), nil
}

func newKv() (*consulapi.KV, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.KV(), nil
}

func newOperator() (*consulapi.Operator, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Operator(), nil
}

func newSession() (*consulapi.Session, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Session(), nil
}

func newStatus() (*consulapi.Status, error) {
	consul, err := newClient()
	if err != nil {
		return nil, err
	}

	return consul.Status(), nil
}

func newClient() (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()

	// Consul address handling

	address := viper.GetString("consul")
	if address != "" {
		config.Address = address
	}

	// ACL token handling

	token := viper.GetString("token")
	tokenFile := viper.GetString("token-file")

	if token != "" && tokenFile != "" {
		return nil, errors.New("--token and --token-file can not both be provided")
	}

	if tokenFile != "" {
		b, err := ioutil.ReadFile(tokenFile)
		if err != nil {
			return nil, fmt.Errorf("Error reading token file: %s", err)
		}
		config.Token = strings.TrimSpace(string(b))
	}

	if token != "" {
		config.Token = token
	}

	// SSL handling

	sslEnabled := viper.GetBool("ssl")
	if sslEnabled {
		config.Scheme = "https"

		tlsConfig := new(tls.Config)
		transport := new(http.Transport)

		sslCert := viper.GetString("ssl-cert")
		sslKey := viper.GetString("ssl-key")
		sslCaCert := viper.GetString("ssl-ca-cert")
		sslVerify := viper.GetBool("ssl-verify")

		if sslCert != "" {
			if sslKey != "" {
				return nil, errors.New("--ssl-key must be provided in order to use certificates for authentication")
			}
			clientCert, err := tls.LoadX509KeyPair(sslCert, sslKey)
			if err != nil {
				return nil, err
			}

			tlsConfig.Certificates = []tls.Certificate{clientCert}
			tlsConfig.BuildNameToCertificate()
		}

		if sslVerify {
			if sslCaCert == "" {
				return nil, errors.New("--ssl-ca-cert must be provided in order to use certificates for verification")
			}

			caCert, err := ioutil.ReadFile(sslCaCert)
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
		config.HttpClient.Transport = transport
	}

	// Auth handling

	authString := viper.GetString("auth")

	if authString != "" {
		auth := new(consulapi.HttpBasicAuth)
		if strings.Contains(authString, ":") {
			split := strings.SplitN(authString, ":", 2)
			auth.Username = split[0]
			auth.Password = split[1]
		} else {
			auth.Username = authString
		}
		config.HttpAuth = auth
	}

	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return client, err
}

func writeOptions() *consulapi.WriteOptions {
	writeOpts := new(consulapi.WriteOptions)

	token := viper.GetString("token")
	if token != "" {
		writeOpts.Token = token
	}

	dc := viper.GetString("datacenter")
	if dc != "" {
		writeOpts.Datacenter = dc
	}

	return writeOpts
}

func queryOptions() *consulapi.QueryOptions {
	queryOpts := new(consulapi.QueryOptions)

	if token := viper.GetString("token"); token != "" {
		queryOpts.Token = token
	}

	if dc := viper.GetString("datacenter"); dc != "" {
		queryOpts.Datacenter = dc
	}

	if wi := viper.Get("wait-index"); wi != nil {
		queryOpts.WaitIndex = wi.(uint64)
	}

	queryOpts.RequireConsistent = viper.GetBool("consistent")
	queryOpts.AllowStale = viper.GetBool("stale")

	return queryOpts
}

func addDatacenterOption(cmd *cobra.Command) {
	cmd.Flags().String("datacenter", "", "Consul data center")
}

func addWaitIndexOption(cmd *cobra.Command) {
	cmd.Flags().Uint64("wait-index", 0, "Only return if ModifyIndex is greater than <index>")
}

func addConsistencyOptions(cmd *cobra.Command) {
	cmd.Flags().Bool("consistent", false, "Enable strong consistency")
	cmd.Flags().Bool("stale", false, "Allow any agent to service the request")
}
