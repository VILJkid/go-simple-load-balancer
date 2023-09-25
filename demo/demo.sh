#!/bin/bash

# Define global variables
server_name="example-server.app"
load_balancer_name="load-balancer.app"
declare -a port_numbers=("8081" "8082")

# Declare an associative array to store Server process IDs
declare -A server_pids

# Declare a variable to store Load Balancer process ID
lb_pid=

# --------------------------------------------------------

# Function to start the servers
start_servers() {
    # Navigate to the example-server directory
    cd example-server
    
    # Perform necessary operations
    rm -rf logs
    mkdir logs
    go mod tidy
    go build -o "$server_name"

    # Start the servers and populate the pids array
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
}

# --------------------------------------------------------

# Function to start the load balancer
start_load_balancer() {
    # Navigate to the load-balancer directory
    cd ../load-balancer

    # Perform necessary operations
    rm -rf logs
    mkdir logs
    go mod tidy
    go build -o "$load_balancer_name"

    # Start the load balancer and populate lb_pid
    ./"$load_balancer_name" "${port_numbers[@]}" > "logs/load-balancer.log" 2>&1 &
    lb_pid=$!
    echo "Load Balancer on port 3000 with pid $lb_pid started."
}

# Function to stop the load balancer gracefully
stop_load_balancer() {
    if [ -n "$lb_pid" ]; then
        if [ -n $safe_exit ]; then
            kill -INT "$lb_pid"
            wait "$lb_pid"
        fi
        echo "Load Balancer on port 3000 with pid $lb_pid stopped."
    fi
}

# --------------------------------------------------------

# Function to start both servers and load balancer
start_all() {
    start_servers
    start_load_balancer
}

# Function to stop both servers and load balancer
stop_all() {
    stop_servers
    stop_load_balancer
    exit 0
}

# --------------------------------------------------------

trap stop_all EXIT
start_all

read -p "Press Enter to stop the servers and load balancer gracefully..."
safe_exit=1