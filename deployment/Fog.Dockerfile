FROM golang:1.22.3 as builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libzmq3-dev \
    libczmq-dev

WORKDIR /app

COPY ./../go.mod ./
COPY ./../go.sum ./

RUN go mod download

ENV GOOS=linux \
    GOARCH=amd64

COPY ../. ./
RUN go build -o /app/fog ./cmd/fog/main.go ./cmd/fog/client.go

FROM debian:12.6-slim as prod

RUN apt-get update && apt-get install -y nmap

WORKDIR /app

COPY --from=builder /app/fog ./fog

COPY ../.env ./.env
COPY ../.fog.env ./.fog.env

EXPOSE 5556
EXPOSE 5557

ENTRYPOINT ["./fog"]
