# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ingredient-detector ./cmd/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 app && adduser -D -u 1000 -G app app

# Set working directory
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/ingredient-detector .

# Create logs directory
RUN mkdir -p logs && chown -R app:app logs

# Expose port
EXPOSE 8080

# Run the application
CMD ["./ingredient-detector"]
