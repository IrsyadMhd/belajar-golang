# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code (including subdirectories)
COPY . .

# Build the application
# CGO_ENABLED=0 for static binary (no C dependencies)
# -ldflags="-s -w" to strip debug info and reduce binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# Production stage - using alpine for smaller image
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy .env file if exists (optional, can also use environment variables)
COPY --from=builder /app/.env* ./

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]
