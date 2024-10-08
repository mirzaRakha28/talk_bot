# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire Go project (including main.go) into the container
COPY . .

# Copy the .env file
COPY .env ./

# Build the Go binary and name it 'seatalk-bot'
RUN go build -o seatalk-bot ./main.go

# Stage 2: Build the final lightweight container
FROM alpine:latest

# Set the working directory for the final container
WORKDIR /root/

# Copy the Go binary from the builder container
COPY --from=builder /app/seatalk-bot .

# Copy the .env file
COPY --from=builder /app/.env ./

# Expose port 5030 for the application
EXPOSE 5030

# Run the Go binary when the container starts
CMD ["./seatalk-bot"]
