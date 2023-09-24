FROM golang:1.20-bullseye AS builder

COPY . /musicFromVideo/

WORKDIR /musicFromVideo/

RUN go mod download
RUN go build -o bin/bot cmd/bot/main.go

FROM buildpack-deps:bullseye-scm

RUN apt-get update && apt-get install curl -y --no-install-recommends && apt-get install ffmpeg  -y --no-install-recommends

RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
RUN chmod a+rx /usr/local/bin/yt-dlp

WORKDIR /root/

COPY --from=builder /musicFromVideo/bin/bot .
COPY --from=builder /musicFromVideo/config config/

EXPOSE 80

CMD ["./bot"]


