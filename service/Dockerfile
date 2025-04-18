# Use the official Golang image to build the application
FROM golang:1.23 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a minimal image for the final container
FROM alpine:latest

# Install libc6-compat for compatibility with Go binaries
RUN apk add --no-cache libc6-compat

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port the Go application will run on
EXPOSE 8080

# Command to run the Go application
CMD ["./main"]