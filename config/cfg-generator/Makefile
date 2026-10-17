.PHONY: build test lint tidy

.DEFAULT_GOAL := build

build:
	go build -o cfg-generator main.go

test:
	go test ./... -v

lint:
	golangci-lint run

tidy:
	go mod tidy