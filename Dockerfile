FROM golang:1.26-bookworm AS bilder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY .. ..

RUN go build -o cryptocurrency ./cmd/cryptocurrency
