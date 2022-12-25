.PHONY: build run build-image start-container

build:
	go build -o ./.bin/bot cmd/main.go

run: build
	./.bin/bot

build-image:
	docker build -t telegram-bot-pocketer:v2.0 .

start-container:
	docker run --env-file .env -p 80:80 telegram-bot-pocketer:v2.0