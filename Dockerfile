# Stage 1: Build
FROM golang:1.24.5 AS builder

# Prebuild Go std
RUN CGO_ENABLED=0 go install std

# Install dlv for debugging
RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 go build -o bot ./cmd/bot

# Stage 2: Run
FROM alpine:3.21 AS production

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app .
COPY --from=builder /go/bin/dlv /

# Expose port (change if needed)
EXPOSE 4000

# Run the application
CMD ["./bot"]
# CMD ["/dlv", "--listen=:4000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./bot"]
