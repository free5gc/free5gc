.PHONY: build test test-run-sh lint tidy

.DEFAULT_GOAL := build

build:
	go build -o cfg-generator main.go

test:
	go test ./... -v

# Runs the repository's run.sh with stub NFs (Linux only).
test-run-sh:
	bash scripts/test-run-sh.sh

lint:
	golangci-lint run

tidy:
	go mod tidy