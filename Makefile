APP=api
PKG=./...


.PHONY: run test build tidy up


run:
	go run ./cmd/api

test:
	go test $(PKG)

build:
	mkdir -p bin
	go build -o bin/$(APP) ./cmd/api

tidy:
	go mod tidy

up:
	podman-compose up -d

down:
	podman-compose down

logs:
	podman-compose logs -f