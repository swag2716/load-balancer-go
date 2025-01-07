package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func handleErr(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type Server interface {
	Address() string
	IsAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type SimpleServer struct {
	Addr  string
	Proxy httputil.ReverseProxy
}

func NewSimpleServer(address string) *SimpleServer {
	serverURL, err := url.Parse(address)
	handleErr(err)
	server := &SimpleServer{
		Addr:  address,
		Proxy: *httputil.NewSingleHostReverseProxy(serverURL),
	}

	return server
}

func InitLoadBalancer(port string, servers []Server) *LoadBalancer {
	LoadBalancer := &LoadBalancer{
		Port:             port,
		RoundRobbinCount: 0,
		Servers:          servers,
	}

	return LoadBalancer
}

func (lb *LoadBalancer) GetNextAvailableServer() Server {
	servers := lb.Servers
	server := servers[lb.RoundRobbinCount%len(servers)]
	for !server.IsAlive() {
		lb.RoundRobbinCount++
		server = servers[lb.RoundRobbinCount%len(servers)]
	}
	lb.RoundRobbinCount++

	return server
}

func (lb *LoadBalancer) ServeProxy(rw http.ResponseWriter, r *http.Request) {
	server := lb.GetNextAvailableServer()
	fmt.Println("forwarding request to:", server.Address())
	server.Serve(rw, r)

}

func (s *SimpleServer) IsAlive() bool {
	return true
}
func (s *SimpleServer) Address() string {
	return s.Addr
}
func (s *SimpleServer) Serve(rw http.ResponseWriter, r *http.Request) {
	s.Proxy.ServeHTTP(rw, r)
}

type LoadBalancer struct {
	Port             string
	RoundRobbinCount int
	Servers          []Server
}

func main() {

	servers := []Server{NewSimpleServer("https://www.google.com/"), NewSimpleServer("http://www.duckduckgo.com/"), NewSimpleServer("https://www.instagram.com/")}

	lb := InitLoadBalancer("8000", servers)

	handleRedirect := func(rw http.ResponseWriter, r *http.Request) {
		lb.ServeProxy(rw, r)
	}

	http.HandleFunc("/", handleRedirect)

	fmt.Println("Listening on port:", lb.Port)

	http.ListenAndServe(":"+lb.Port, nil)
}
