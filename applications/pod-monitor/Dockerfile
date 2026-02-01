# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o watcher .

# Final stage
FROM alpine:latest

# Add ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Copy only the binary from builder
COPY --from=builder /app/watcher .

# Change ownership
RUN chown -R appuser:appgroup /app

# Use non-root user
USER appuser

CMD ["./watcher"]
