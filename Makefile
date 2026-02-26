APP=api
PKG=./...
MIGRATIONS_DIR=db/migrations

.PHONY: run test build tidy up down logs migrate-up migrate-down migrate-version migrate-create


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

migrate-up:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up
migrate-down:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1
migrate-version:
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version
migrate-create:
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
