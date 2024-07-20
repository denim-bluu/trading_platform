# deployments/Dockerfile
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Build the Go app
ARG SERVICE_NAME
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/service ./cmd/${SERVICE_NAME}

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/service .

# Copy the .env file
COPY deployments/.env .

# Install the 'env' command
RUN apk add --no-cache gettext

# Expose port 50051 to the outside world
EXPOSE 50051

# Use env to load environment variables and run the service
CMD ["/bin/sh", "-c", "env $(cat .env | xargs) ./service"]