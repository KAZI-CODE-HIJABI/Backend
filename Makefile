.PHONY: run build test check migrate-up migrate-down up down
run:
	go run ./cmd/api
build:
	go build -o bin/api ./cmd/api
test:
	go test ./...
check:
	go vet ./...
	go test ./...
migrate-up:
	go run ./cmd/migrate up
migrate-down:
	go run ./cmd/migrate down
up:
	docker compose up --build -d
down:
	docker compose down
