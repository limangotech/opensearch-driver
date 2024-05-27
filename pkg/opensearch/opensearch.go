//go:generate mockgen -build_flags=--mod=mod -destination ../../tests/mocks/mock_opensearchapi/transport_mock.go github.com/opensearch-project/opensearch-go/v2/opensearchapi Transport

package opensearch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/golang-migrate/migrate/v4/database"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"go.uber.org/atomic"
)

const (
	nullVersion                     = -1
	versionIndexName                = ".migrations"
	errTemplateUnsupportedOperation = "unsupported operation '%s'"
)

type OpenSearch struct {
	transport         opensearchapi.Transport
	manager           MigrationsIndexManagerInterface
	MigrationSequence []string
	LastRunMigration  []byte
	isLocked          atomic.Bool
}

func NewDriver(
	transport opensearchapi.Transport,
	manager MigrationsIndexManagerInterface,
) *OpenSearch {
	return &OpenSearch{
		transport:         transport,
		manager:           manager,
		MigrationSequence: make([]string, 0),
	}
}

func (o *OpenSearch) Lock() error {
	if !o.isLocked.CompareAndSwap(false, true) {
		return database.ErrLocked
	}

	return nil
}

func (o *OpenSearch) Unlock() error {
	if !o.isLocked.CompareAndSwap(true, false) {
		return database.ErrNotLocked
	}

	return nil
}

func (o *OpenSearch) Run(migration io.Reader) error {
	content, err := io.ReadAll(migration)
	if err != nil {
		return fmt.Errorf("could not read migration content on run: %w", err)
	}

	mgr, err := NewMigrationFromRawContent(content)
	if err != nil {
		return err
	}

	req, err := mgr.CreateRequest()
	if err != nil {
		return err
	}

	resp, err := o.transport.Perform(req)
	defer closeResponseBody(resp)

	if err != nil {
		return fmt.Errorf("apply migration request failed: %w", err)
	}

	if err = ReadErrorFromResponse(resp); err != nil {
		return err
	}

	o.LastRunMigration = content
	o.MigrationSequence = append(o.MigrationSequence, string(content))

	return nil
}

func (o *OpenSearch) SetVersion(version int, dirty bool) error {
	mgr := Migration{
		Method: "PUT",
		URL:    fmt.Sprintf("/%s/_doc/1", versionIndexName),
		Body: map[string]any{
			"version": version,
			"dirty":   dirty,
		},
	}
	req, err := mgr.CreateRequest()

	if err != nil {
		return err
	}

	resp, err := o.transport.Perform(req)
	defer closeResponseBody(resp)

	if err != nil {
		return fmt.Errorf("set version request failed: %w", err)
	}

	return ReadErrorFromResponse(resp)
}

//nolint:nonamedreturns
func (o *OpenSearch) Version() (version int, dirty bool, err error) {
	if err = o.manager.Upsert(versionIndexName); err != nil {
		return nullVersion, false, fmt.Errorf("could not upsert migrations index: %w", err)
	}

	m := Migration{
		Method: "GET",
		URL:    fmt.Sprintf("/%s/_doc/1", versionIndexName),
	}
	req, err := m.CreateRequest()

	if err != nil {
		return nullVersion, false, err
	}

	resp, err := o.transport.Perform(req)
	defer closeResponseBody(resp)

	if err != nil {
		return nullVersion, false, fmt.Errorf("get version request failed: %w", err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return nullVersion, false, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nullVersion, false, fmt.Errorf("could not read body from response on getting version: %w", err)
	}

	var parsedResp getDocumentResponse

	if err = json.Unmarshal(body, &parsedResp); err != nil {
		return nullVersion, false, fmt.Errorf("could not unmarshal response body JSON on getting version: %w", err)
	}

	return parsedResp.Source.Version, parsedResp.Source.Dirty, nil
}

//nolint:ireturn
func (o *OpenSearch) Open(_ string) (database.Driver, error) {
	return nil, fmt.Errorf(errTemplateUnsupportedOperation, "open")
}

func (o *OpenSearch) Close() error {
	return fmt.Errorf(errTemplateUnsupportedOperation, "close")
}

func (o *OpenSearch) Drop() error {
	return fmt.Errorf(errTemplateUnsupportedOperation, "drop")
}
