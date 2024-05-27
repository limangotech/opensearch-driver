.PHONY: all init generate test lint

all: lint test

init:
	go run -mod=mod github.com/google/wire/cmd/wire ./...

generate:
	rm -f **/*_gen.go **/*_mock.go
	go generate ./...

lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

test:
	go test -v ./...