.PHONY:

build-image:
	docker build -t telegram-bot:v0.1 .

start-container:
	docker run --name telegram-bot -p 80:80 --env-file .env telegram-bot:v0.1 --restart always