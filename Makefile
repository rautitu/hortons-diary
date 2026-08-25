.PHONY: build test fmt compose-up compose-down db-up run

build:
	go build ./cmd/server

test:
	go test ./...

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

compose-up:
	docker compose up --build

compose-down:
	docker compose down

db-up:
	docker compose -f compose.yaml -f compose.dev-db.yaml up -d db

run:
	go run ./cmd/server
