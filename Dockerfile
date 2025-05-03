FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o snell-panel .

# Create a minimal image for running the application
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite

# Create a non-root user to run the application
RUN adduser -D -h /app appuser
USER appuser
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder --chown=appuser:appuser /app/snell-panel /app/

# Create directory for SQLite database with proper permissions
RUN mkdir -p /app/data && chown -R appuser:appuser /app/data

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./snell-panel"]