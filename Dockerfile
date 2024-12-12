# Build stage
FROM golang:1.23.2-alpine AS builder

# Set working directory
WORKDIR /build

# Install required system dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-w -s" -o main .

# Final stage
FROM alpine:3.19

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/main .
COPY --from=builder /build/migrations ./migrations
COPY --from=builder /build/views ./views

# Expose port 8080
EXPOSE 1323

# Command to run the executable
CMD ["./main"]
