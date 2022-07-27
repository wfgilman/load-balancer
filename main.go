package main

import (
	"fmt"

	"github.com/wfgilman/balancer/backend"
	"github.com/wfgilman/balancer/client"
)

func main() {
	servers := []backend.WebServer{
		*backend.NewWebServer("https://www.google.com"),
		*backend.NewWebServer("https://www.amazon.com"),
	}

	lb := &client.LoadBalancer{
		Port:    8000,
		Current: 0,
	}

	for _, server := range servers {
		lb.AddBackend(&server)
	}
}
