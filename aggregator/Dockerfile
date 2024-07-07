# Start from the official Go image
FROM golang:1.22.5-alpine

# Set the working directory inside the container
WORKDIR /app/trading_platform

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Set the working directory to the aggregator folder
WORKDIR /app/trading_platform/aggregator

# Build the application
RUN go build -o main ./cmd/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]