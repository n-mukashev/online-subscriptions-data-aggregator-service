# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o subscriptions-app ./cmd/main.go

# Stage 2: Final image
FROM alpine:latest
WORKDIR /app

# Устанавливаем ca-certificates (нужно для HTTPS в golang-migrate)
RUN apk add --no-cache ca-certificates

# Копируем бинарь, конфиг и миграции
COPY --from=builder /app/subscriptions-app .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./subscriptions-app"]
