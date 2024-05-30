# Use the official Golang image to create a build artifact.
FROM golang:1.22.3 as builder

# Create and change to the working directory.
WORKDIR /app

# Copy the go.mod file from the root directory.
COPY ./../go.mod ./

# Retrieve application dependencies.
RUN go mod download

# Copy the source code from the cmd/fog directory and build the fog binary.
COPY ../cmd/fog ./cmd/fog
RUN CGO_ENABLED=0 GOOS=linux go build -v -o fog-server ./cmd/fog

# Copy the source code from the cmd/cloud directory and build the cloud binary.
COPY ../cmd/cloud ./cmd/cloud
RUN CGO_ENABLED=0 GOOS=linux go build -v -o cloud-server ./cmd/cloud

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

# Copy the binaries to the production image from the builder stage.
COPY --from=builder /app/fog-server /fog-server
COPY --from=builder /app/cloud-server /cloud-server

