package backend

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createWebServer(t *testing.T, addr string) *WebServer {
	serverUrl, _ := url.Parse(addr)
	return &WebServer{
		Addr:  addr,
		Url:   serverUrl,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

func TestAddress(t *testing.T) {
	addr := "http://localhost:3000"
	webServer := createWebServer(t, addr)

	got := webServer.Address()

	assert.Equal(t, addr, got)
}

func TestIsAlive(t *testing.T) {
	webServer := createWebServer(t, "http://localhost:3000")

	expect := true
	got := webServer.IsAlive()

	assert.Equal(t, expect, got)
}

func TestServeHTTP(t *testing.T) {
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This would be one of our many backend servers rendering
		// responses to the client.
	}))
	defer backend.Close()

	backendUrl, err := url.Parse(backend.URL)
	if err != nil {
		t.Fatal(err)
	}

	webServer := createWebServer(t, backendUrl.String())
	reverseProxy := httptest.NewServer(webServer) // This is the line that calls the ServerHTTP method on the WebServer struct.
	defer reverseProxy.Close()

	res, err := http.Get(reverseProxy.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	expect := fmt.Sprintf("Hello from server at %s\n", webServer.Addr)
	actual := string(body)

	assert.Equal(t, expect, actual)
}

func TestNewWebServer(t *testing.T) {
	addr := "http://localhost:3000"
	expect := createWebServer(t, addr)
	actual := NewWebServer(addr)

	assert.Equal(t, expect.Addr, actual.Addr)
	assert.Equal(t, expect.Alive, actual.Alive)
}
