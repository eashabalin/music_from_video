.PHONY:

build-image:
	docker build -t eashabalin/music_from_youtube .

start-container:
	docker run --name telegram-bot -p 80:80 --env-file .env --restart=always -d eashabalin/music_from_youtube