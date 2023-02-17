FROM golang:alpine AS builder

RUN apk update && apk add ffmpeg && apk add youtube-dl

COPY . /musicFromVideo/

WORKDIR /musicFromVideo/

RUN go mod download
RUN go build -o bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /musicFromVideo/bin/bot .
COPY --from=builder /musicFromVideo/config config/

EXPOSE 80

CMD ["./bot"]


