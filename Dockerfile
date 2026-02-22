# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod file
COPY go.mod ./

# Download dependencies (if any)
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# CGO_ENABLED=0 ensures a static binary, crucial for scratch or alpine images
RUN CGO_ENABLED=0 GOOS=linux go build -o http-debug .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/http-debug .

# Expose the port the app runs on
EXPOSE 8080

# Run the binary
CMD ["./http-debug"]
