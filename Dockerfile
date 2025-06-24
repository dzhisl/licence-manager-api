# --- Build Stage ---
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o license-api ./cmd/server/main.go

# --- Runtime Stage ---
FROM gcr.io/distroless/static-debian11
WORKDIR /app
COPY --from=builder /app/license-api .
# Optionally copy Swagger docs if you want to serve them statically
COPY docs ./docs
COPY .env .
EXPOSE 8080
CMD ["./license-api"]
