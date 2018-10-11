package commands

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckRegister(t *testing.T) {
	// wait channel for the http check
	requestChan := make(chan struct{})
	tested := false
	fakeService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "testhost", r.Host)
		assert.True(t, assert.ObjectsAreEqual([]string{"TestVal1,TestVal2"}, r.Header["X-Test-Header"]))
		if !tested {
			tested = true
			requestChan <- struct{}{}
		}
	}))
	defer fakeService.Close()

	// register the service with `consul-cli`
	serviceArgs := []string{"service", "register", "test-check-service", "--consul", consulTestAddr}
	checkArgs := []string{"check", "register", "test-http-check", "--interval", "1s", "--http", fakeService.URL, "--header", "Host: testhost", "--header", "X-Test-Header: TestVal1,TestVal2"}

	command := NewConsulCliCommand("consul-cli", "0.0.1")
	command.SetArgs(serviceArgs)
	err := command.Execute()
	assert.Nil(t, err)

	command.ResetFlags()

	command.SetArgs(checkArgs)
	err = command.Execute()
	assert.Nil(t, err)

	select {
	case <-requestChan:
		return
	case <-time.After(2 * time.Second):
		t.Fatalf("Timeout waiting for health check")
	}
}
