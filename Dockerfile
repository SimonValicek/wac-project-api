# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./
RUN go mod download

# Copy everything else
COPY . .

# Build the binary
RUN go build -o wac-api-service ./cmd/wac-api-service

# Final stage - smaller image with just the binary and necessary files
FROM alpine:latest

# Optional: add CA certs for HTTPS if needed
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy the built binary from builder stage
COPY --from=builder /app/wac-api-service .

# Expose port 8080
EXPOSE 8080

# Database environment (can override in k8s)
ENV AMBULANCE_API_ENVIRONMENT=production 
ENV AMBULANCE_API_PORT=8080 
ENV AMBULANCE_API_MONGODB_HOST=mongo 
ENV AMBULANCE_API_MONGODB_PORT=27017 
ENV AMBULANCE_API_MONGODB_DATABASE=xvaliceks-wac-project 
ENV AMBULANCE_API_MONGODB_COLLECTION=parking_lots 
ENV AMBULANCE_API_MONGODB_USERNAME=root 
ENV AMBULANCE_API_MONGODB_PASSWORD= 
ENV AMBULANCE_API_MONGODB_TIMEOUT_SECONDS=5

# Run the binary
CMD ["./wac-api-service"]