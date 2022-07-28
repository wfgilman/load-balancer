package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/wfgilman/balancer/backend"
	"github.com/wfgilman/balancer/client"
)

func main() {
	serverList := flag.String("backends", "localhost:3000", "Enter backends in format 'localhost:3000,localhost:3001'")
	port := flag.Int("port", 8000, "Enter port number for load balancer")
	tokens := strings.Split(*serverList, ",")

	var servers []backend.WebServer

	for _, token := range tokens {
		servers = append(servers, *backend.NewWebServer(token))
	}

	lb := &client.LoadBalancer{
		Port:    *port,
		Current: 0,
	}

	for _, server := range servers {
		go server.StartServer()
		lb.AddBackend(&server)
	}

	lbHandler := func(rw http.ResponseWriter, req *http.Request) {
		lb.ServeProxy(rw, req)
	}

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", lb.Port),
		Handler: http.HandlerFunc(lbHandler),
	}

	log.Printf("Load Balancer started at :%d\n", lb.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
