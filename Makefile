.PHONY: build run

build:
	docker build -t siphon .

run:
	docker run --rm -it siphon
