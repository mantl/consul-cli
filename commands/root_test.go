package commands

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/consul/sdk/testutil"
)

var consulTestAddr string

func TestMain(m *testing.M) {
	var consulServerOutputEnabled bool
	flag.BoolVar(&consulServerOutputEnabled, "enable-consul-output", false, "Enables consul server output")
	flag.Parse()

	// create a test Consul server
	consulTestServer, err := testutil.NewTestServerConfig(func(c *testutil.TestServerConfig) {
		c.Connect = nil

		if !consulServerOutputEnabled {
			c.Stderr = ioutil.Discard
			c.Stdout = ioutil.Discard
		}
	})

	if err != nil {
		panic(err)
	}

	consulTestAddr = consulTestServer.HTTPAddr

	exitCode := m.Run()
	consulTestServer.Stop()

	os.Exit(exitCode)
}
