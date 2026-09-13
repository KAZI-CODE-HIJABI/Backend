.PHONY: run build test check up down
run:
	go run ./cmd/api
build:
	go build -o bin/api ./cmd/api
test:
	go test ./...
check:
	go vet ./...
	go test ./...
up:
	docker compose up --build -d
down:
	docker compose down
