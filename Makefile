# Makefile for go-template

REPO:=gbaeke
IMAGE:=$(REPO)/go-template

test:
	go test -v -race ./...

build:
	CGO_ENABLED=0 go build -installsuffix 'static' -o app cmd/app/*

docker:
	docker build -t $(IMAGE) .