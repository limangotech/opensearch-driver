package opensearch_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/dhojayev/opensearch-driver/pkg/opensearch"
	"github.com/dhojayev/opensearch-driver/tests/mocks/mock_opensearchapi"
	"github.com/dhojayev/opensearch-driver/tests/stubs"
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

	// @case create if does not exist
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
	if err != nil {
		t.Error(err)
	}

	if exists != true {
		t.Error("Expected true, got false")
	}

	// @case returns false on 404
	resp = http.Response{StatusCode: http.StatusNotFound}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	exists, err = manager.Exists("test")
	if err != nil {
		t.Error(err)
	}

	if exists != false {
		t.Error("Expected false, got true")
	}

	// @case returns false on error
	transport.EXPECT().Perform(gomock.Any()).Return(nil, errors.New("test error"))

	exists, err = manager.Exists("test")
	if err == nil {
		t.Error("Expected error, got nil")
	}

	if exists != false {
		t.Errorf("Expected false, got true")
	}
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
	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	// @case returns error on failure
	resp = http.Response{
		StatusCode: http.StatusBadGateway,
		Body:       &stubs.ResponseBody{Reader: strings.NewReader("test error")},
	}

	transport.EXPECT().Perform(gomock.Any()).Return(&resp, nil)

	err = manager.Create("test-index")
	if err == nil {
		t.Error("Expected error, got nil")

		return
	}

	if !strings.Contains(err.Error(), "test error") {
		t.Errorf("Unexpected error: %s", err)
	}
}
