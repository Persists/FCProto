# Use the official Golang image to create a build artifact.
FROM golang:1.22.3 as builder

# Install necessary dependencies.
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libzmq3-dev \
    libczmq-dev

# Create and change to the working directory.
WORKDIR /app

# Copy the go.mod and go.sum files from the root directory.
COPY ./../go.mod ./
COPY ./../go.sum ./

# Retrieve application dependencies.
RUN go mod download

# Set environment variables.
ENV GOOS=linux \
    GOARCH=arm64

# Copy the source code and build the fog binary.
COPY ../cmd/fog ./cmd/fog
RUN go build -v -o fog-server ./cmd/fog

COPY ../cmd/cloud ./cmd/cloud
RUN go build -v -o cloud-server ./cmd/cloud
