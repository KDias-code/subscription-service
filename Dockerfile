FROM golang:1.25.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscriptions-service ./cmd/main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache bash curl

COPY --from=builder /app/subscriptions-service .
COPY ./manifests ./manifests

EXPOSE 8080

CMD ["./subscriptions-service"]
