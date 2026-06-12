.PHONY: build run test seed clean

build:
	go build -o bin/api ./cmd/api
	go build -o bin/seeder ./cmd/seeder

run:
	go run ./cmd/api

test:
	go test -v ./...

seed:
	go run ./cmd/seeder

clean:
	rm -rf bin/ hadith.db