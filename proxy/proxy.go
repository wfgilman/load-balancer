package proxy

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/wfgilman/balancer/backend"
	"github.com/wfgilman/balancer/pool"
	"github.com/wfgilman/balancer/utils"
)

func New(targetAddr string) *httputil.ReverseProxy {
	targetUrl := &url.URL{
		Scheme: "http",
		Host:   targetAddr,
	}
	return httputil.NewSingleHostReverseProxy(targetUrl)
}

func ErrorHandler(p *httputil.ReverseProxy, sp pool.Pool, be *backend.Backend) func(rw http.ResponseWriter, req *http.Request, e error) {
	// A custom error handler which will retry a failed request then
	// server another backend if the request continues to fail.
	return func(rw http.ResponseWriter, req *http.Request, e error) {
		log.Printf("[%s] %s\n", be.Address(), e.Error())
		retries := utils.CountRetries(req)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond):
				ctx := context.WithValue(req.Context(), utils.Retries, retries+1)
				p.ServeHTTP(rw, req.WithContext(ctx))
			}
			return
		}

		// After 3 retries to the same server, mark as down.
		server, err := sp.GetServer(be.Address())
		if err != nil {
			log.Printf("[%s] %s\n", be.Address(), err)
		}
		server.SetAlive(false)

		// Increment the number of attempts of the request and route it to a
		// new server using the pool algorithm.
		attempts := utils.CountAttempts(req)
		log.Printf("[%s] Attemping retry %d\n", req.RemoteAddr, attempts)
		ctx := context.WithValue(req.Context(), utils.Attempts, attempts+1)
		sp.Serve(rw, req.WithContext(ctx))
	}
}
