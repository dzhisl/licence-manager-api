# --- Build Stage ---
FROM golang:1.24-alpine AS builder

# Установим необходимые пакеты для сборки
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o license-api ./cmd/server/main.go

# --- Runtime Stage ---
FROM alpine:latest

# add shell for docker exec
RUN apk add --no-cache bash

WORKDIR /app

COPY --from=builder /app/license-api .
COPY docs ./docs
COPY .env .

EXPOSE 8080

CMD ["./license-api"]
