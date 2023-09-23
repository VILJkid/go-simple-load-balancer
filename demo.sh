#!/bin/bash

# Navigate to the example-server directory
cd example-server
    
# Perform necessary Go build operations
app_name="example-server.bin"
go mod tidy
go build -o "$app_name"

# Define an array of port numbers
declare -a port_numbers=("8081" "8082")

# Declare an associative array to store pids
declare -A pids

# Function to start the servers
start_servers() {
    # Start the Go servers and populate the pids array
    for port in "${port_numbers[@]}"; do
        ./"$app_name" "$port" > "server$port.log" 2>&1 &
        pids["$port"]=$!
        echo "Server on port $port with pid ${pids[$port]} started."
    done
}

# Function to stop the servers gracefully
stop_servers() {
    echo ""
    for port in "${port_numbers[@]}"; do
        pid="${pids[$port]}"
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