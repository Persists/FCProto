FROM golang:1.22.3 as builder

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libzmq3-dev \
    libczmq-dev

WORKDIR /app

COPY ../../go.mod ./
COPY ../../go.sum ./

RUN go mod download

ENV GOOS=linux \
    GOARCH=amd64

COPY ../.. ./
RUN go build -o /app/edge ./cmd/edge/main.go

FROM debian:12.6-slim as prod

RUN apt-get update && apt-get install -y nmap

WORKDIR /app

COPY --from=builder /app/edge ./edge

EXPOSE 5556

ENTRYPOINT ["./edge"]
