// main.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"log/slog"

	"github.com/VILJkid/go-simple-load-balancer/example-server/server"
)

const usage = `Usage:
	go build -o example-server
	./example-server <port>`

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	port := os.Args[1]

	// Create and set a logger for better logging
	slog.SetDefault(slog.Default())

	// Initialize the HTTP server
	server := server.NewServer(port)

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server running on :", "port", port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error:", "error", err)
		}
	}()

	// Block until a signal is received for graceful shutdown
	<-stop

	// Perform any necessary cleanup and shutdown operations in the server
	server.Shutdown()
}
