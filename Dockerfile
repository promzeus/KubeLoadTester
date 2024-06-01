# syntax=docker/dockerfile:1

FROM golang:1.21-bookworm AS build
WORKDIR /src

# Копируем go.mod и go.sum файлы и загружаем зависимости
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod tidy && go mod download

# Копируем остальную часть исходного кода и строим проект
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go build -o /out/bin

# Используем более легковесный образ для финального контейнера
FROM debian:bookworm
STOPSIGNAL SIGTERM
ENTRYPOINT [ "bash", "/docker-entrypoint.sh" ]
CMD ["start"]

# Устанавливаем необходимые зависимости
RUN apt update && apt install -y ca-certificates

# Копируем скрипт запуска
COPY docker-entrypoint.sh /

# Копируем собранное бинарное приложение из стадии сборки
COPY --from=build /out/bin /go/run

# Открываем порт, если это необходимо
EXPOSE 8080