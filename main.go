package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wfgilman/balancer/backend"
	"github.com/wfgilman/balancer/pool"
	"github.com/wfgilman/balancer/proxy"
)

func healthCheck() {
	t := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-t.C:
			log.Println("Health check started")
			serverPool.HealthCheck()
			log.Println("Health check complete")
		}
	}
}

func serverStats() {
	t := time.NewTicker(20 * time.Minute)
	for {
		select {
		case <-t.C:
			log.Println("Server Report...")
			serverPool.ServerStats()
		}
	}
}

var serverPool pool.Pool

func main() {
	// Parse startup parameters from command line.
	numBackends := flag.Int("n", 5, "Enter number of backend servers")
	port := flag.Int("port", 8000, "Enter port number for load balancer")
	algo := flag.String("algo", "alwaysfirst", "Balancing algorithm")
	flag.Parse()

	serverPool.Algorithm = *algo

	for i := 0; i < *numBackends; i++ {
		// Create backend server.
		addr := fmt.Sprintf("localhost:500%d", i)
		be := backend.New(addr)
		beServer := http.Server{
			Addr:    addr,
			Handler: http.HandlerFunc(be.RequestHandler()),
		}

		// Start backend server.
		log.Printf("(%s) Backend Server started\n", be.Address())
		go beServer.ListenAndServe()

		// Randomly crash a server
		if i%*numBackends == 1 {
			time.AfterFunc(25*time.Second, func() {
				log.Printf("(%s) Backend Server crashed\n", be.Address())
				beServer.Shutdown(context.Background())
			})
		}

		// Create a ReverseProxy pointing to backend server.
		p := proxy.New(be.Address())
		p.ErrorHandler = proxy.ErrorHandler(p, serverPool, be)
		be.SetReverseProxy(p)

		// Add the backend to the server pool.
		serverPool.AddServer(be)
	}

	// Create the load balancer server.
	server := http.Server{
		Addr:    fmt.Sprintf("localhost:%d", *port),
		Handler: http.HandlerFunc(serverPool.RequestHandler()),
	}

	go healthCheck()
	go serverStats()

	// Start the load balancer.
	log.Printf("(%s) Load Balancer started\n", server.Addr)
	server.ListenAndServe()
}
