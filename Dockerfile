# Используем базовый образ для Go
FROM golang:latest

RUN go version

ADD app .

RUN go mod download



RUN GOOS=linux go build -o ./.bin/bot ./cmd/bot/main.go

WORKDIR /root/

EXPOSE 80

CMD ["./bot"]