FROM golang:1.23-alpine AS builder

# Set environment variables for Go module and x64 architecture
ARG ARCH=amd64
ARG GO111MODULE=on

RUN apk add --no-cache git ca-certificates

WORKDIR /code

# Copy go mod and sum files and download all the dependencies
COPY go.mod go.sum ./
COPY . /code

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} go build \
    -o /go/bin/file-storage \
    -gcflags "all=-trimpath=/go" \
    -a ./cmd/file-storage

# Run the container
FROM alpine:3.15

RUN apk --no-cache add ca-certificates wget

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /go/bin/file-storage /bin/file-storage

# Copy config file
COPY --from=builder /code/cmd/file-storage/config.yaml ./config.yaml

EXPOSE 8080

ENTRYPOINT ["/bin/file-storage", "-config-path", "/app/config.yaml"]

