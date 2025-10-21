#!/bin/bash

services=("file-storage" "user-management" "ui")

generate_mock() {
    service=$1
    type=$2
    path=$3
    dest=$4

    echo "Generating mock for $service $type..."
    docker run -t --rm -v "$(pwd)":/code krelms/golang-dev-tools mockgen -source=/code/"$path" -destination /code/"$dest"
    echo "Completed mock generation for $service $type"
}

for service in "${services[@]}"; do
    generate_mock "$service" "service" "pkg/$service/services/interface.go" "pkg/$service/services/mocks/mock.go" &
    generate_mock "$service" "repository" "pkg/$service/repository/interface.go" "pkg/$service/repository/mocks/mock.go" &
done

generate_mock "queue" "interface" "pkg/queue/interface.go" "pkg/queue/mocks/mock.go" &
wait

echo "All mock generations completed."