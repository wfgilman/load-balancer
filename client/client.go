package client

import (
	"fmt"
	"net/http"

	"github.com/wfgilman/balancer/backend"
)

type LoadBalancer struct {
	Port     int
	Current  uint64
	Backends []*backend.WebServer
}

func (lb *LoadBalancer) AddBackend(backend *backend.WebServer) {
	lb.Backends = append(lb.Backends, backend)
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request) {
	proxyServer := lb.Backends[0]
	fmt.Printf("Forwarding request to address %s\n", proxyServer.Address())
	proxyServer.ServeHTTP(rw, req)
}
