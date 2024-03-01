package opensearch

import (
	"fmt"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type MigrationsIndexManagerInterface interface {
	Upsert(name string) error
	Exists(name string) (bool, error)
	Create(name string) error
}

type MigrationsIndexManager struct {
	transport opensearchapi.Transport
}

func NewMigrationsIndexManager(transport opensearchapi.Transport) MigrationsIndexManager {
	return MigrationsIndexManager{
		transport: transport,
	}
}

func (m MigrationsIndexManager) Upsert(name string) error {
	exists, err := m.Exists(name)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	return m.Create(name)
}

func (m MigrationsIndexManager) Exists(name string) (bool, error) {
	mgr := Migration{
		Method: "HEAD",
		URL:    "/" + name,
	}
	req, err := mgr.CreateRequest()

	if err != nil {
		return false, err
	}

	resp, err := m.transport.Perform(req)
	defer closeResponseBody(resp)

	if err != nil {
		return false, fmt.Errorf("check migrations index Exists request failed: %w", err)
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (m MigrationsIndexManager) Create(name string) error {
	mgr := Migration{
		Method: "PUT",
		URL:    "/" + name,
		Body: map[string]any{
			"mappings": map[string]any{
				"properties": map[string]any{
					"version": map[string]any{
						"type": "integer",
					},
					"dirty": map[string]any{
						"type": "boolean",
					},
				},
			},
		},
	}
	req, err := mgr.CreateRequest()

	if err != nil {
		return err
	}

	resp, err := m.transport.Perform(req)
	defer closeResponseBody(resp)

	if err != nil {
		return fmt.Errorf("create migrations index request failed: %w", err)
	}

	return ReadErrorFromResponse(resp)
}
