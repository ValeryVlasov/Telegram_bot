FROM golang:1.19-alpine3.17 AS builder

COPY . /github.com/ValeryVlasov/Telegram_bot/
WORKDIR /github.com/ValeryVlasov/Telegram_bot/

RUN go mod download
RUN go build -o ./bin/bot cmd/bot/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /github.com/ValeryVlasov/Telegram_bot/bin/bot .
COPY --from=0 /github.com/ValeryVlasov/Telegram_bot/configs configs/

EXPOSE 80

CMD ["./bot"]