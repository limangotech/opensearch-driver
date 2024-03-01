.PHONY: test-lint

mocks:
	go install go.uber.org/mock/mockgen@latest
	mockgen -destination tests/mocks/mock_opensearchapi/transport.go github.com/opensearch-project/opensearch-go/v2/opensearchapi Transport

test-unit:
	go test ./...

test-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2
	golangci-lint run ./...