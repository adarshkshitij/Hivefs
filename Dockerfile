# Stage 1: Build the binary
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures the binary is statically linked
RUN CGO_ENABLED=0 GOOS=linux go build -o /fs-node .

# Stage 2: Run the binary
FROM alpine:latest

# Add non-root user for security
RUN adduser -D fsuser
USER fsuser

WORKDIR /home/fsuser

# Copy the binary from the builder stage
COPY --from=builder /fs-node .

# Default environment variables
ENV LISTEN_ADDR=":3000"
ENV STORAGE_ROOT="data"
ENV METRICS_ADDR=":9090"

# Create storage directory
RUN mkdir -p data

# Expose the application and metrics ports
EXPOSE 3000 9090

# Run the binary
CMD ["./fs-node"]
