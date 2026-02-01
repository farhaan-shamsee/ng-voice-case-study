# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go module files first (for better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY main.go .

# Build statically-linked binary with maximum optimization
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o pod-watcher .

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/pod-watcher /pod-watcher

USER nonroot:nonroot

ENTRYPOINT ["/pod-watcher"]
