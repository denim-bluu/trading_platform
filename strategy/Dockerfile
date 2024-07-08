FROM golang:1.22.5-alpine

WORKDIR /app

# Set the working directory inside the container
WORKDIR /app/trading_platform

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project into the container
COPY . .

# Set the working directory to the aggregator folder
WORKDIR /app/trading_platform/strategy

# Build the application
RUN go build -o strategy ./cmd/main.go

EXPOSE 8081

CMD ["./strategy"]
