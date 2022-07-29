package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/wfgilman/balancer/backend"
	"github.com/wfgilman/balancer/client"
)

func main() {
	// Parse startup parameters from command line.
	serverList := flag.String("backends", "localhost:3000", "Enter backends in format: localhost:4000,localhost:4001")
	port := flag.Int("port", 8000, "Enter port number for load balancer")
	flag.Parse()

	// Setup context and channel for graceful server shutdown.
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-c
		cancel()
	}()

	// Create Load Balancer.
	lb := &client.LoadBalancer{
		Port:    *port,
		Current: 0,
	}

	// Start backend servers in a group of go routines
	// so we can send a shutdown message to the group.
	grp, grpContext := errgroup.WithContext(ctx)
	tokens := strings.Split(*serverList, ",")
	for _, token := range tokens {
		webServer := backend.NewWebServer(token)

		server := http.Server{
			Addr:    webServer.Address(),
			Handler: webServer,
		}

		lb.AddBackend(webServer)

		grp.Go(func() error {
			log.Printf("Backend WebServer started at %s\n", webServer.Address())
			return server.ListenAndServe()
		})
		grp.Go(func() error {
			<-grpContext.Done()
			log.Printf("Backend WebServer at %s shutdown gracefully\n", webServer.Address())
			return server.Shutdown(context.Background())
		})
	}

	lbHandler := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}

	balancer := http.Server{
		Addr:    fmt.Sprintf(":%d", lb.Port),
		Handler: http.HandlerFunc(lbHandler),
	}

	// Run the Load Balancer in the same group so we can shut everything down
	// gracefully with Ctrl+C from the command line.
	grp.Go(func() error {
		log.Printf("Load Balancer started at :%d\n", lb.Port)
		return balancer.ListenAndServe()
	})
	grp.Go(func() error {
		<-grpContext.Done()
		log.Printf("Load Balancer shutdown gracefully")
		return balancer.Shutdown(context.Background())
	})

	if err := grp.Wait(); err != nil {
		fmt.Printf("Program exited with reason: %s\n", err)
	}
}
