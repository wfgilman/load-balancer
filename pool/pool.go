package pool

import (
	"errors"
	"log"
	"math"
	"net"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/wfgilman/balancer/server"
)

const (
	AlwaysFirst  string = "alwaysfirst"
	RoundRobin          = "roundrobin"
	LeastLatency        = "leastlatency"
	FewestConn          = "fewestconn"
)

type Pool struct {
	Current   uint64
	Servers   []server.Server
	Algorithm string
}

func (p *Pool) AddServer(server server.Server) {
	p.Servers = append(p.Servers, server)
}

func (p *Pool) RequestHandler() func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		p.Serve(rw, req)
	}
}

func (p *Pool) Serve(rw http.ResponseWriter, req *http.Request) {
	server := p.GetNextServer()
	proxy := server.ReverseProxy()
	proxy.ServeHTTP(rw, req)
}

func (p *Pool) GetNextServer() server.Server {
	switch p.Algorithm {
	case AlwaysFirst:
		return p.Servers[0]
	case RoundRobin:
		nextIndex := int(atomic.AddUint64(&p.Current, uint64(1)) % uint64(len(p.Servers)))
		length := len(p.Servers) + nextIndex
		for i := nextIndex; i < length; i++ {
			index := i % len(p.Servers)
			if p.Servers[index].IsAlive() {
				if i != nextIndex {
					atomic.StoreUint64(&p.Current, uint64(index))
				}
				return p.Servers[index]
			}
		}
		panic("No healthy backends exist")
	case LeastLatency:
		var s server.Server
		min := math.MaxInt
		for _, server := range p.Servers {
			latency := server.AverageLatency()
			if server.IsAlive() && latency < min {
				min = latency
				s = server
			}
		}
		return s
	case FewestConn:
		var s server.Server
		min := math.MaxInt
		for _, server := range p.Servers {
			activeConn := server.ActiveConnections()
			if server.IsAlive() && activeConn < min {
				min = activeConn
				s = server
			}
		}
		return s
	}
	panic("No backends exist")
}

func (p *Pool) GetServer(targetAddr string) (server.Server, error) {
	for _, server := range p.Servers {
		if server.Address() == targetAddr {
			return server, nil
		}
	}
	return nil, errors.New("Server not found")
}

func (p *Pool) HealthCheck() {
	for _, server := range p.Servers {
		if server.IsAlive() {
			alive := isServerAlive(server)
			server.SetAlive(alive)
			if !alive {
				log.Printf("(%s) is down\n", server.Address())
			}
		}
	}
}

func isServerAlive(server server.Server) bool {
	url := &url.URL{
		Scheme: "http",
		Host:   server.Address(),
	}
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", url.Host, timeout)
	if err != nil {
		log.Printf("[%s] unreachable, error: ", server.Address(), err)
		return false
	}
	defer conn.Close()
	return true
}

func (p *Pool) ServerStats() {
	for _, s := range p.Servers {
		log.Printf("(%s) Alive %v, Active %d, Total %d, Latency %d(ms)\n", s.Address(), s.IsAlive(), s.ActiveConnections(), s.TotalRequests(), s.AverageLatency())
	}
}
