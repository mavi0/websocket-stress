# Stage 1: Build the application
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
# Statically link the binary to avoid C library dependencies
# Specify the output path for the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/websocket-stress ./cmd/websocket-stress

# Stage 2: Create the final lightweight image
FROM alpine:latest

WORKDIR /app

# Copy static files and templates
COPY --from=builder /app/web /app/web

# Copy the built binary from the builder stage
COPY --from=builder /app/websocket-stress /app/websocket-stress

# Expose the port the application runs on
EXPOSE 8088

# Command to run the application
CMD ["/app/websocket-stress"] 