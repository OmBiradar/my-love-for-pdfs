# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the application
ARG SERVICE
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-bin ./cmd/${SERVICE}

# Run stage - using a smaller Alpine image
FROM alpine:latest

# Install Ghostscript which is required for PDF compression
RUN apk add --no-cache ghostscript

# Create a non-root user to run the application
RUN adduser -D appuser
USER appuser

# Copy the binary from the build stage
COPY --from=builder /app-bin /app-bin

# Set proper permissions
USER root
RUN chmod +x /app-bin
USER appuser

# Create temp directory with proper permissions for temporary files
RUN mkdir -p /tmp/pdf-files
WORKDIR /tmp/pdf-files

# Define port exposure and entrypoint
EXPOSE 8080
CMD ["/app-bin"]