# Build stage
FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the CDC application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o cdc ./cmd/cdc

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/cdc .

# Run the CDC application
CMD ["./cdc", "--rabbitmq-host=rabbitmq", "--typesense-host=typesense"]
