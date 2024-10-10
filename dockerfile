# Use the official Go image as a base image
FROM golang:1.23-alpine

# Set environment variables
ENV GO111MODULE=on

# Set the working directory in the container
WORKDIR /app

# Copy go.mod and go.sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of your application code
COPY . .

# Accept build arguments
ARG SEATALK_APP_ID
ARG SEATALK_APP_SECRET
ARG REGRESSION_GROUP_ID
ARG SEATALK_API_URL
ARG SEATALK_AUTH_URL
ARG SEATALK_SEND_SINGLE_CHAT_URL
ARG SEATALK_SEND_GROUP_CHAT_URL
ARG PORT

# Set environment variables using the build arguments
ENV SEATALK_APP_ID=${SEATALK_APP_ID}
ENV SEATALK_APP_SECRET=${SEATALK_APP_SECRET}
ENV REGRESSION_GROUP_ID=${REGRESSION_GROUP_ID}
ENV SEATALK_API_URL=${SEATALK_API_URL}
ENV SEATALK_AUTH_URL=${SEATALK_AUTH_URL}
ENV SEATALK_SEND_SINGLE_CHAT_URL=${SEATALK_SEND_SINGLE_CHAT_URL}
ENV SEATALK_SEND_GROUP_CHAT_URL=${SEATALK_SEND_GROUP_CHAT_URL}
ENV PORT=${PORT}

# Build the Go application
RUN go build -o main .

# Expose the port your application runs on
EXPOSE 6969

# Command to run the executable
CMD ["./main"]
