package client

import (
	"fmt"
	"net/http"
	"sync/atomic"

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

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, req *http.Request, webServer *backend.WebServer) {
	fmt.Printf("Forwarding request to address %s\n", webServer.Address())
	webServer.Serve(rw, req)
}

func (lb *LoadBalancer) AlwaysFirst() *backend.WebServer {
	return lb.Backends[0]
}

func (lb *LoadBalancer) RoundRobin() *backend.WebServer {
	nextIndex := int(atomic.AddUint64(&lb.Current, uint64(1)) % uint64(len(lb.Backends)))
	return lb.Backends[nextIndex]
}
