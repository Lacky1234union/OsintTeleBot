# Build stage
FROM golang:1.22.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy any additional required files (like configs, migrations, etc.)
COPY --from=builder /app/migrations ./migrations

# Set timezone
ENV TZ=UTC

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"] 