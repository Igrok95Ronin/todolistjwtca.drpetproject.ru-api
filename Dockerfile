FROM golang:1.24.0 AS build

WORKDIR /usr/local/src

# Копируем файлы зависимостей и загружаем их
COPY go.mod go.sum ./

# Копируем весь проект
COPY . .

# Проверяем версию Go (опционально, для диагностики)
RUN go version

# Сборка приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /todolistjwtca-drpetproject ./cmd/todolistjwtca

# Этап финального образа
FROM alpine:3.18 AS runner

# Устанавливаем bash, если необходимо
RUN apk --no-cache add bash

WORKDIR /usr/local/src

# Копируем скомпилированный бинарник и конфигурационный файл
COPY --from=build /todolistjwtca-drpetproject /usr/local/src/todolistjwtca-drpetproject
COPY ./config.yml ./

# Определяем точку входа
ENTRYPOINT ["/usr/local/src/todolistjwt-drpetproject"]