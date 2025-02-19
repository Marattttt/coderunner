#!/bin/bash
set -oue

# Check if an argument is provided
if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <cli|rest>"
    exit 1
fi

MODE=$1

if [[ "$MODE" == "cli" ]]; then
    echo "Running in CLI mode..."
    docker build . -f Dockerfile.cli -t runnercli
    docker run --cap-add=NET_ADMIN --env-file .env runnercli:latest
elif [[ "$MODE" == "rest" ]]; then
    echo "Running in REST mode..."
    docker build . -f Dockerfile.rest -t runnerrest
    docker run -p 8080:8080 --cap-add=NET_ADMIN --env-file .env runnerrest:latest
else
    echo "Invalid argument: $MODE"
    echo "Usage: $0 <cli|rest>"
    exit 1
fi

