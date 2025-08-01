# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies for CGO
RUN apk add --no-cache build-base

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY *.go ./

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o wsrelay .

# Runtime stage
FROM alpine:latest

# Install required dependencies for SQLite
RUN apk add --no-cache libc6-compat

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/wsrelay .

# Create data directory for SQLite database
RUN mkdir -p /app/data

# Expose the port the app runs on
EXPOSE 8080

# Create a volume for persistent data storage
VOLUME ["/app/data"]

# Run the application
CMD ["./wsrelay"]