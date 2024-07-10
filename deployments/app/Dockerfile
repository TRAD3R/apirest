# Этап 1: Сборка
FROM golang:1.22.5-alpine3.20 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем остальные файлы проекта
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 go build -o /app/bin/apirest ./cmd/app/main.go

# Этап 2: Запуск
FROM alpine:3.20.1

WORKDIR /app

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/bin/apirest .

# Устанавливаем команду для запуска приложения
ENTRYPOINT ["/app/apirest"]
