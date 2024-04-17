.PHONY: all install generate test lint

all: lint test

install:
	go install go.uber.org/mock/mockgen@latest
	export PATH=$PATH:$(go env GOPATH)/bin

generate:
	go generate -v ./...

test:
	go test -v ./...

lint:
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2 run -v ./...