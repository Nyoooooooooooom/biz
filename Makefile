APP=biz
BIN=./bin/$(APP)

.PHONY: tidy build test test-offline run-local clean docker-build

tidy:
	go mod tidy

build:
	go build -o $(BIN) ./cmd/biz

test:
	go test ./...

test-offline:
	mkdir -p .gocache .gomodcache
	GOCACHE=$(PWD)/.gocache GOMODCACHE=$(PWD)/.gomodcache go test ./internal/invoice ./internal/invoice/datasource ./internal/invoice/notion ./internal/platform/errors ./internal/platform/output ./internal/platform/clock ./internal/tax

run-local:
	go run ./cmd/biz --config ./config.example.yaml invoice list --status ready --json

clean:
	rm -rf ./bin ./.gocache ./.gomodcache

docker-build:
	docker build -t biz:latest .
