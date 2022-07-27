package client

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wfgilman/balancer/backend"
)

func TestAddBackend(t *testing.T) {
	backendUrl, _ := url.Parse("http://localhost:3000")
	webServer := &backend.WebServer{
		Addr:  "http://localhost:3000",
		Url:   backendUrl,
		Alive: true,
	}

	lb := &LoadBalancer{
		Port:    3030,
		Current: 0,
	}

	lb.AddBackend(webServer)

	assert.Equal(t, lb.Backends[0], webServer)
}

func TestServeProxy(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This would be one of our many backend servers rendering
		// responses to the client.
	}))
	defer server.Close()

	serverUrl, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	webServer := &backend.WebServer{
		Addr:  serverUrl.String(),
		Url:   serverUrl,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}

	lb := LoadBalancer{
		port: 8000,
	}
}
