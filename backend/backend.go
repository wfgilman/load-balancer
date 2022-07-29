package backend

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, req *http.Request)
}

type WebServer struct {
	Addr  string
	Url   *url.URL
	Alive bool
	Proxy *httputil.ReverseProxy
}

func (ws *WebServer) Address() string {
	return ws.Addr
}

func (ws *WebServer) IsAlive() bool {
	return ws.Alive
}

func (ws *WebServer) Serve(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Hello from server at %s\n", ws.Addr)
	ws.Proxy.ServeHTTP(rw, req)
}

func NewWebServer(addr string) *WebServer {
	serverUrl, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}

	return &WebServer{
		Addr:  addr,
		Url:   serverUrl,
		Alive: true,
		Proxy: httputil.NewSingleHostReverseProxy(serverUrl),
	}
}
