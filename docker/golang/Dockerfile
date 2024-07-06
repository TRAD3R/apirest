# Этап 1: Сборка
FROM golang:latest AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum для установки зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем остальные файлы проекта
COPY . .

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/myapp ./cmd/app/main.go

# Этап 2: Запуск
FROM alpine:latest

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /root/

# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/bin/myapp .

# Устанавливаем переменную окружения для использования без libc
ENV CGO_ENABLED=0

# Устанавливаем команду для запуска приложения
CMD ["./myapp"]
