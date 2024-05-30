# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.22.3 as builder

# Create and change to the working directory.
WORKDIR /app

# Copy the go.mod file to the working directory.
COPY go.mod ./

# Copy the source code to the working directory.
COPY src/fog ./

# Build the binary.
# -o specifies the output file
# Using CGO_ENABLED=0 to ensure a fully static binary that can run without C dependencies.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o server

# Use a Docker multi-stage build to create a lean production image.
# Start with a smaller image without extra layers.
FROM alpine:latest

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Run the web service on container startup.
CMD ["/server"]