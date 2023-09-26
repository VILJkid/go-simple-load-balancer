package server

import (
	"math/rand"
	"net/url"
	"os"
	"sync"
	"sync/atomic"
)

type Server struct {
	URL         *url.URL
	ActiveConns int32
}

// List of backend servers
var servers []*Server

func NewServer(port string) *Server {
	return &Server{
		URL: &url.URL{
			Scheme: "http",
			Host:   ":" + port,
		},
		ActiveConns: 0,
	}
}

func SetBackendServers() {
	serverPorts := os.Args[1:]
	for _, port := range serverPorts {
		servers = append(servers, NewServer(port))
	}
}

var serverLock sync.Mutex

// Assume weights are stored in an array `weights`
// such that weights[i] corresponds to servers[i]
var weights = []int{3, 1} // example weights

// Load balancing logic: Random
func GetRandomBackend() *Server {
	index := rand.Intn(len(servers))
	return servers[index]
}

// Load balancing logic: Least Connections
func GetLeastConnectionsBackend() *Server {
	serverLock.Lock()
	defer serverLock.Unlock()

	// Find the backend with the least connections
	var leastConnServer *Server
	for _, server := range servers {
		if leastConnServer == nil || atomic.LoadInt32(&server.ActiveConns) < atomic.LoadInt32(&leastConnServer.ActiveConns) {
			leastConnServer = server
		}
	}

	// Increment the active connection count for this server
	atomic.AddInt32(&leastConnServer.ActiveConns, 1)

	return leastConnServer
}

// Load balancing logic: Weighted Random
func GetWeightedRandomBackend() *Server {
	totalWeight := 0
	for _, weight := range weights {
		totalWeight += weight
	}

	randomValue := rand.Intn(totalWeight)

	for i, weight := range weights {
		randomValue -= weight
		if randomValue < 0 {
			return servers[i]
		}
	}

	return servers[len(servers)-1] // fallback, should not happen
}

// Load balancing logic: Weighted Random Algorithm combined with Least Connections
func GetDynamicWeightedRandomBackend() *Server {
	// Calculate the total dynamic weight
	totalWeight := int32(0)
	for _, server := range servers {
		// Assume the weight is inversely proportional to the number of active connections
		// Add 1 to avoid division by zero
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		totalWeight += int32(weight)
	}

	// Select a backend server based on dynamic weight
	randomValue := int32(rand.Intn(int(totalWeight)))
	for _, server := range servers {
		weight := 1 / (atomic.LoadInt32(&server.ActiveConns) + 1)
		randomValue -= int32(weight)
		if randomValue < 0 {
			atomic.AddInt32(&server.ActiveConns, 1) // Increment active connections
			return server
		}
	}

	return nil // fallback, should not happen
}

// Decrement connection count for a server
func ReleaseBackend(server *Server) {
	atomic.AddInt32(&server.ActiveConns, -1)
}
