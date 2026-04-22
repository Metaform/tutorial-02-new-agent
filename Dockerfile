# Build stage
FROM --platform=$BUILDPLATFORM golang:1.26.2 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

# Cache dependencies separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -ldflags="-s -w" -o server cmd/server/main.go

# Runtime stage
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Run the binary
ENTRYPOINT ["./server"]
