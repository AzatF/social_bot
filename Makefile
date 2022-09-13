.PHONY: default build image run

default: build

build:
	CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' cmd/main.go

image:
	docker build -t soc-bot-v3:latest .
	docker image prune -f --filter label=stage=tgbotbuilder2

run:
	go run cmd/main.go