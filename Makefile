.PHONY: build run test seed clean docs build-mcp run-mcp

build:
	go build -o bin/api ./cmd/api
	go build -o bin/seeder ./cmd/seeder

build-mcp:
	go build -o bin/mcp-server ./cmd/mcp

run:
	go run ./cmd/api

run-mcp:
	go run ./cmd/mcp

test:
	go test -v ./...

seed:
	go run ./cmd/seeder

docs:
	swag init -g cmd/api/main.go -o docs/swagger

clean:
	rm -rf bin/ hadith.db