# Use the official Go image as the base image
FROM golang:1.17-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files to the container
COPY go.mod go.sum ./

# Download and cache Go modules dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go application
RUN go build -o blog-api .

# Use a lightweight base image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/blog-api .

# Expose the desired port (change if needed)
EXPOSE 8080

# Set the command to run the binary
CMD ["./blog-api"]
