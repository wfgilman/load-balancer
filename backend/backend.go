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
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
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

func (ws *WebServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
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

func (ws *WebServer) StartServer() {
	server := http.Server{
		Addr:    ws.Address(),
		Handler: ws,
	}

	log.Printf("Backend WebServer started at %s\n", ws.Address())
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	defer server.Close()
}
