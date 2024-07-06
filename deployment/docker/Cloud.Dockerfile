FROM golang:1.22.3 as builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libzmq3-dev \
    libczmq-dev

WORKDIR /app

<<<<<<< HEAD
COPY ../../go.mod ./
COPY ../../go.sum ./

=======
COPY ./go.mod ./
COPY ./go.sum ./
>>>>>>> 3293c89 (fix: docker deployment)
RUN go mod download

ENV GOOS=linux \
    GOARCH=amd64

COPY ../.. ./
RUN go build -o /app/cloud ./cmd/cloud/main.go ./cmd/cloud/server.go

FROM debian:12.6-slim as prod

RUN apt-get update && apt-get install -y nmap

WORKDIR /app

COPY --from=builder /app/cloud ./cloud

EXPOSE 5555

ENTRYPOINT ["./cloud"]
