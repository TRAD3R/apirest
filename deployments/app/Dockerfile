# Этап 1: Сборка
FROM golang:1.22.4-alpine3.19 AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем остальные файлы проекта
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 go build -o /app/bin/myapp ./cmd/app/main.go

# Этап 2: Запуск
FROM alpine:3.19

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/bin/myapp .

# Устанавливаем команду для запуска приложения
ENTRYPOINT ["./myapp"]
