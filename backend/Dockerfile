# Stage 1: Build
FROM golang:1.23 as builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code including migrations
COPY . .

# Build the Go app
RUN go build -o social-network server.go
# Stage 2: Runtime container
FROM ubuntu:22.04

WORKDIR /app/backend

RUN apt-get update && \
    apt-get install -y ca-certificates sqlite3 && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copy compiled binary
COPY --from=builder /app/social-network .

# Copy migrations folder from builder stage
COPY --from=builder /app/pkg ./pkg

# Copy static folders (e.g., imageUpload)
COPY --from=builder /app/imageUpload ./imageUpload

RUN mkdir -p ./imageUpload

EXPOSE 8081

CMD ["./social-network"]

