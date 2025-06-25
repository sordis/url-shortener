# Build stage
FROM golang:1.24.4-alpine3.22 AS builder

# Устанавливаем зависимости для CGO и инструменты
RUN apk add --no-cache \
    gcc \
    musl-dev \
    git \
    make

WORKDIR /app
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o url-shortener ./cmd/url-shortener

# Runtime stage
FROM alpine:3.22

# Минимальные зависимости
RUN apk add --no-cache \
    libc6-compat \
    tzdata

WORKDIR /app

# Подготовка папки для БД
RUN mkdir -p /app/storage && \
    chmod -R 777 /app/storage

# Копируем бинарник и конфиги
COPY --from=builder /app/url-shortener .
COPY --from=builder /app/config/prod.yml ./config/


RUN touch ./config/local.yml

# Аргументы сборки
ARG AUTH_PASS

# Переменные окружения
ENV TZ=Europe/Moscow \
    CONFIG_PATH=/app/config/prod.yml \
    AUTH_PASS=$AUTH_PASS

EXPOSE 8080

# Запуск
ENTRYPOINT ["./url-shortener"]