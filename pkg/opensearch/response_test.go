package opensearch_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/limangotech/opensearch-driver/pkg/opensearch"
	"github.com/limangotech/opensearch-driver/tests/stubs"
)

func TestReadErrorFromResponse(t *testing.T) {
	t.Parallel()

	// @case returns nil on success code
	resp := http.Response{
		StatusCode: http.StatusOK,
	}

	err := opensearch.ReadErrorFromResponse(&resp)

	assert.NoError(t, err)

	// @case returns error from body
	resp = http.Response{
		StatusCode: http.StatusBadRequest,
		Body: &stubs.ResponseBody{
			Reader: strings.NewReader("test error"),
		},
	}

	err = opensearch.ReadErrorFromResponse(&resp)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test error")
}
