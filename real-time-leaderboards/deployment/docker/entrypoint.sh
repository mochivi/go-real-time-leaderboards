#!/bin/bash

# Default mode is dev if TEAMSERVER_MODE is not set
MODE=${SERVER_MODE:-dev}
echo "Running in $MODE mode"

# Run migrations
# make migrate-up
# ./migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -verbose up


if [ "$MODE" = "debug" ]; then
    echo ""
else
    # Install https://github.com/air-verse/air for hot reloading
    echo "Installing air..."
    go install github.com/air-verse/air@latest
    echo "Finished installing air"

    echo "Starting server..."
    exec air
fi
