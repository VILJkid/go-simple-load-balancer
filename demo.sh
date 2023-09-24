#!/bin/bash

# Define global variables
server_name="example-server.app"
declare -a port_numbers=("8081" "8082")

# --------------------------------------------------------

# Navigate to the example-server directory
cd example-server
    
# Perform necessary operations
rm -rf logs
mkdir logs
go mod tidy
go build -o "$server_name"

# Define an associative array to store Server process IDs
declare -A server_pids

# Function to start the servers
start_servers() {
    # Start the Go servers and populate the pids array
    for port in "${port_numbers[@]}"; do
        ./"$server_name" "$port" > "logs/server$port.log" 2>&1 &
        server_pids["$port"]=$!
        echo "Server on port $port with pid ${server_pids[$port]} started."
    done
}

# Function to stop the servers gracefully
stop_servers() {
    echo ""
    for port in "${port_numbers[@]}"; do
        pid="${server_pids[$port]}"
        if [ -n "$pid" ]; then
            if [ -n "$safe_exit" ]; then
                kill -INT "$pid"
                wait "$pid"
            fi
            echo "Server on port $port with pid $pid stopped."
        fi
    done
    exit 0
}

# Set up an exit trap to stop servers on script exit
trap stop_servers EXIT

# Start the servers
start_servers

read -p "Press Enter to stop the servers gracefully..."
safe_exit=1

# --------------------------------------------------------