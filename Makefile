# Makefile for go-template

REPO:=gbaeke
TAG:=latest
IMAGE:=$(REPO)/go-template:$(TAG)


test:
	go test -v -race ./...

build:
	CGO_ENABLED=0 go build -installsuffix 'static' -o app cmd/app/*

docker-build:
	docker build -t $(IMAGE) .

docker-push:
	docker build -t $(IMAGE) .
	docker push $(IMAGE)

swagger:
	cd pkg/api ; swag init -g server.go