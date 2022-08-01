# Balancer
A simple load balancer written in Go.

### Features
* [ ] Testing
* [x] Round Robin Method
* [x] Health Checks
* [x] Fewest Active Connections Method
* [x] Least Latency Method
* [ ] Weighted Round Robin

### How it works
![Schematic](/assets/lb.png)

On start-up, the application creates `n` backend web servers in the `5000` block on
the `localhost` port from `5000` to `500N`. Each backend conforms to the `Server` interface.
The backend superficially responds to `GET` requests with its address and the response time.
Response times are randomly chosen between 0 and 250ms.

Each backend is assigned a single host reverse proxy. The purpose of the proxy in this
application is to handle backend failure gracefully within the server pool. In Go, the
reverse proxy has an `ErrorHandler` method which can be customized to retry requests or
propagate issues with its assigned backend to the server pool.

Each backend is assigned to the server pool of the load balancer. The pool is responsible
for implementing the server rotation algorithms. An HTTP server is created and implements
a response handler for the pool that decides which backend to send the request to.

The `main` function runs a health check on an interval which makes a TCP connection to each backend
server to verify it is alive, then sets the status of each backend through the `Server` interface
accordingly.

### Round Robin Algorithm
![Schematic](/assets/go-balancer.gif)

### Credits
The design of this application is inspired and informed by the content and code
in the online resources below.

* Load Balancer: https://kasvith.me/posts/lets-create-a-simple-lb-go/
* Load Balancer: https://betterprogramming.pub/building-a-load-balancer-in-go-3da3c7c46f30
* Graceful Server Shutdown: https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go
* Reverse Proxy: https://blog.joshsoftware.com/2021/05/25/simple-and-powerful-reverseproxy-in-go/
* Enums: https://www.sohamkamani.com/golang/enums/
