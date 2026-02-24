APP=api
PKG=./...


.PHONY: run test build tidy


run:
	go run ./cmd/api

test:
	go test $(PKG)

build:
	mkdir -p bin
	go build -o bin/$(APP) ./cmd/api

tidy:
	go mod tidy