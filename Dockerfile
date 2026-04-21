# Build stage
FROM golang:1.26.2 AS builder

WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Run the binary
ENTRYPOINT ["./server"]
