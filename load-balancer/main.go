package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/VILJkid/go-simple-load-balancer/balancer/server"
)

// var servers = []*Server{
// 	{
// 		URL: &url.URL{
// 			Scheme: "http",
// 			Host:   ":8081",
// 		},
// 		ActiveConns: 0,
// 	},
// 	{
// 		URL: &url.URL{
// 			Scheme: "http",
// 			Host:   ":8082",
// 		},
// 		ActiveConns: 0,
// 	},
// }

// HTTP handler
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	backend := server.GetDynamicWeightedRandomBackend()
	proxy := httputil.NewSingleHostReverseProxy(backend.URL)
	proxy.ServeHTTP(w, r)
	server.ReleaseBackend(backend)
}

func main() {
	if len(os.Args) < 2 {
		panic("atleast one port is required")
	}

	server.SetBackendServers()

	http.HandleFunc("/", proxyHandler)

	fmt.Println("Load Balancer running on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
}
