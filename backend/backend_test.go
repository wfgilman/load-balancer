package backend

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createBackend(t *testing.T, addr string) *Backend {
	url, _ := url.Parse(addr)
	return &Backend{
		URL: url,
	}
}

func TestNew(t *testing.T) {
	addr := "http://localhost:4000"
	be := createBackend(t, addr)
	backend := New(addr)

	assert.Equal(t, backend.URL, be.URL)
}

func TestAddress(t *testing.T) {
	addr := "https://www.google.com"
	backend := New(addr)

	assert.Equal(t, backend.Address(), addr)
}

func TestIsAlive(t *testing.T) {
	backend := New("http://localhost:4000")

	assert.Equal(t, true, backend.IsAlive())
}

func TestSetAlive(t *testing.T) {
	backend := New("http://localhost:4000")
	backend.SetAlive(false)

	assert.Equal(t, false, backend.IsAlive())
}

func TestServe(t *testing.T) {
	backend := New("http://localhost:4000")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		backend.Serve(rw, req)
	})

	handler.ServeHTTP(rr, req)

	assert.NotEqual(t, rr.Body.String(), "")
	assert.Equal(t, backend.ActiveConnections(), 0)
	if backend.AverageLatency() <= 0 {
		t.Error("Expected average latency > 0")
	}
}
