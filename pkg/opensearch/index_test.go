package opensearch_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/limangotech/opensearch-driver/pkg/opensearch"
	"github.com/limangotech/opensearch-driver/tests/mocks/mock_opensearchapi"
	"github.com/limangotech/opensearch-driver/tests/stubs"
)

func TestMigrationsIndexManager_Upsert(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	transport := mock_opensearchapi.NewMockTransport(ctrl)
	manager := opensearch.NewMigrationsIndexManager(transport)

	// @case return on exist
	resp := http.Response{StatusCode: http.StatusOK}

	transport.
		EXPECT().
		Perform(gomock.Any()).
		Return(&resp, nil).
		Times(1)

	_ = manager.Upsert("test index")

	// @case create if it does not exist
	resp = http.Response{StatusCode: http.StatusNotFound}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil).Times(1)
	transport.EXPECT().Perform(gomock.Any()).Return(&http.Response{}, nil).Times(1)

	_ = manager.Upsert("test index")
}

func TestMigrationsIndexManager_Exists(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	transport := mock_opensearchapi.NewMockTransport(ctrl)
	manager := opensearch.NewMigrationsIndexManager(transport)

	// @case returns true on 200
	resp := http.Response{StatusCode: http.StatusOK}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	exists, err := manager.Exists("test")

	assert.NoError(t, err)
	assert.True(t, exists)

	// @case returns false on 404
	resp = http.Response{StatusCode: http.StatusNotFound}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	exists, err = manager.Exists("test")

	assert.NoError(t, err)
	assert.False(t, exists)

	// @case returns false on error
	transport.EXPECT().Perform(gomock.Any()).Return(nil, errors.New("test error"))

	exists, err = manager.Exists("test")

	assert.Error(t, err)
	assert.False(t, exists)
}

func TestMigrationsIndexManager_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	transport := mock_opensearchapi.NewMockTransport(ctrl)
	manager := opensearch.NewMigrationsIndexManager(transport)

	// @case returns nil on success
	resp := http.Response{StatusCode: http.StatusOK}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	err := manager.Create("test-index")

	assert.NoError(t, err)

	// @case returns error on failure
	resp = http.Response{
		StatusCode: http.StatusBadGateway,
		Body:       &stubs.ResponseBody{Reader: strings.NewReader("test error")},
	}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	err = manager.Create("test-index")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test error")
}
