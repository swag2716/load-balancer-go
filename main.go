package main

import(
	"fmt"
	"net/http/httputil"
	"het/http"
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWrite, r *http.Request)
}

type simpleServer struct{
	address string
	proxy *httputil.ReverseProxy

}

type LoadBalancer struct {
	port string
	roundRobbinCount int
	servers []Server
}

func newSimpleServer(add string) *simpleServer{
	serverUrl, err := url.Parse(add) 
	handleErr(err)

	return &simpleServer{
		address: add,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl)
	}
}

func handleErr(err error){
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

func main(){
	
}