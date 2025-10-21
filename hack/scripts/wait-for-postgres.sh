#!/bin/bash

CONTAINER_NAME=$1

echo "Waiting for PostgreSQL container '$CONTAINER_NAME' to be ready..."

max_attempts=10
attempt=0

while [ $attempt -lt $max_attempts ]; do
    if docker exec $CONTAINER_NAME pg_isready -U postgres > /dev/null 2>&1; then
        echo "PostgreSQL container '$CONTAINER_NAME' is ready!"
        exit 0
    fi

    attempt=$((attempt + 1))
    echo "Attempt $attempt/$max_attempts: PostgreSQL not ready yet..."
    sleep 2
done

echo "ERROR: PostgreSQL container '$CONTAINER_NAME' did not become ready in time"
exit 1

