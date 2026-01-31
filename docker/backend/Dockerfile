# Build stage
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Копируем go.mod и go.sum
COPY backend/go.mod backend/go.sum ./
# ДОБАВЛЕН ПРЕФИКС backend/

# Скачиваем зависимости
RUN go mod download

# Копируем весь код бэкенда
COPY backend/ ./
# ДОБАВЛЕН ПРЕФИКС backend/

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Собираем миграции
RUN CGO_ENABLED=0 GOOS=linux go build -o migrate ./migrate.go

# Запуск воркеров
RUN CGO_ENABLED=0 GOOS=linux go build -o articles ./cmd/worker/articles.go
RUN CGO_ENABLED=0 GOOS=linux go build -o stat ./cmd/worker/stat.go

# Production stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/articles .
COPY --from=builder /app/stat .

EXPOSE 8080

CMD ["sh", "-c", "./migrate && ./main"]