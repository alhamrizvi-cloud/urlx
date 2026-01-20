# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o urlx main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 urlx && \
    adduser -D -u 1000 -G urlx urlx

# Set working directory
WORKDIR /home/urlx

# Copy binary from builder
COPY --from=builder /app/urlx .

# Change ownership
RUN chown -R urlx:urlx /home/urlx

# Switch to non-root user
USER urlx

# Set entrypoint
ENTRYPOINT ["./urlx"]

# Default command (show help)
CMD ["-h"]
