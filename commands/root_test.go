package commands

import (
	"os"
	"testing"

	"github.com/hashicorp/consul/testutil"
)

var consulTestAddr string

func TestMain(m *testing.M) {
	// create a test Consul server
	consulTestServer, err := testutil.NewTestServerConfig(func(c *testutil.TestServerConfig) {
		c.Connect = nil
	})

	if err != nil {
		panic(err)
	}

	consulTestAddr = consulTestServer.HTTPAddr

	exitCode := m.Run()
	consulTestServer.Stop()

	os.Exit(exitCode)
}
