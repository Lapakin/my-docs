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
    -o /go/bin/ui \
    -gcflags "all=-trimpath=/go" \
    -a ./cmd/ui

# Run the container
FROM alpine:3.15
RUN apk --no-cache add ca-certificates wget
WORKDIR /app

# Copy the binary from builder

COPY --from=builder /go/bin/ui /bin/ui
# Copy static files and templates

COPY --from=builder /code/pkg/ui/static ./static
COPY --from=builder /code/pkg/ui/templates ./templates

# Copy config file
COPY --from=builder /code/cmd/ui/config.yaml ./config.yaml
EXPOSE 8082

ENTRYPOINT ["/bin/ui", "-config-path", "/app/config.yaml"]
