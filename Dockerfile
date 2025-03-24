# Step 1: Use an official Go image to build the application
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the Go application
RUN go build -o server .

# Step 2: Use a smaller image to run the application
FROM alpine:latest

# Install necessary dependencies (like SSL certificates, if needed)
RUN apk --no-cache add ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the pre-built binary from the builder container
COPY --from=builder /app/server .

# Expose the port the app will run on
EXPOSE 8080

# Command to run the application
CMD ["./server"]
