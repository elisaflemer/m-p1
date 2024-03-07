#!/bin/sh

# Function to cleanup and stop all processes
cleanup() {
    echo "Stopping all processes..."
    # Terminate the Docker Compose services
    docker-compose down
    # Terminate the Python API
    pkill -f "python3 api.py"
    # Terminate the Go subscriber
    pkill -f "go run subscriber.go"
    # Terminate the Go publisher
    pkill -f "go run publisher.go"
    echo "All processes stopped."
    exit 0
}

# Trap SIGINT signal (Ctrl+C) to execute cleanup function
trap 'cleanup' INT

# Activate the virtual environment
source env/bin/activate

# Start Docker Compose services in detached mode
docker-compose up -d

# Start Python API in the background
python3 api.py &

# Start Go subscriber in the background
go run subscriber.go -config config.json -connection hivemq -username elisa -password Elisa123 &

# Start Go publisher in the background
go run publisher.go -config config.json -csv data.csv -connection hivemq -username elisa -password Elisa123

# Wait for all background processes to finish
wait
