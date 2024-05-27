package opensearch_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/dhojayev/opensearch-driver/pkg/opensearch"
	"github.com/dhojayev/opensearch-driver/tests/stubs"
)

func TestReadErrorFromResponse(t *testing.T) {
	t.Parallel()

	// @case returns nil on success code
	resp := http.Response{
		StatusCode: http.StatusOK,
	}

	err := opensearch.ReadErrorFromResponse(&resp)
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	// @case returns error from body
	resp = http.Response{
		StatusCode: http.StatusBadRequest,
		Body: &stubs.ResponseBody{
			Reader: strings.NewReader("test error"),
		},
	}

	err = opensearch.ReadErrorFromResponse(&resp)
	if err == nil {
		t.Error("Expected error, got nil")

		return
	}

	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("Unexpected error: %s", err)
	}
}
