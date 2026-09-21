FROM golang:1.26-alpine AS bilder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o cryptocurrency ./cmd/cryptocurrency

FROM alpine:3.22

WORKDIR /app

COPY --from=bilder /app/cryptocurrency .
COPY --from=bilder /app/deploy ./deploy

EXPOSE 8080

CMD ["./cryptocurrency"]