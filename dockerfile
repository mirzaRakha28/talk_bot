# Use the official Go image as a base image
FROM golang:1.23-alpine

# Set environment variables
ENV GO111MODULE=on

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod  ./

# Download dependencies
RUN go mod download

# Copy the rest of your application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your application runs on
EXPOSE 7070

# Command to run the executable
CMD ["./main"]
