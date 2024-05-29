package opensearch_test

import (
	"io"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/limangotech/opensearch-driver/pkg/opensearch"
)

func TestNewMigrationFromRawContent(t *testing.T) {
	t.Parallel()

	// @case error on invalid JSON
	rawContent := []byte("invalid")
	_, err := opensearch.NewMigrationFromRawContent(rawContent)

	assert.ErrorContains(t, err, "JSON syntax error")

	// @case error on validation error
	rawContent = []byte(`{
    "method": "",
    "url": "/_index_template/test"
}`)

	_, err = opensearch.NewMigrationFromRawContent(rawContent)

	assert.ErrorContains(t, err, "migration content validation failure")

	// @case valid
	rawContent = []byte(`{
    "method": "DELETE",
    "url": "/_index_template/test"
}`)

	expected := opensearch.Migration{
		Method: "DELETE",
		URL:    "/_index_template/test",
	}
	migration, err := opensearch.NewMigrationFromRawContent(rawContent)

	assert.NoError(t, err)
	assert.Equal(t, expected, migration)
}

func TestMigration_Validate(t *testing.T) {
	t.Parallel()

	migration := opensearch.Migration{
		URL: "/_index_template/test",
	}

	err := migration.Validate()

	assert.ErrorIs(t, err, opensearch.ErrInvalidContent)

	migration = opensearch.Migration{
		Method: "DELETE",
	}

	err = migration.Validate()

	assert.ErrorIs(t, err, opensearch.ErrInvalidContent)

	migration = opensearch.Migration{
		Method: "DELETE",
		URL:    "/_index_template/test",
	}

	err = migration.Validate()

	assert.NoError(t, err)
}

func TestMigration_CreateRequest(t *testing.T) {
	t.Parallel()

	migration := opensearch.Migration{
		Method: "DELETE",
		URL:    "/_index_template/test",
	}

	request, err := migration.CreateRequest()

	assert.NoError(t, err)
	assert.Equal(t, "DELETE", request.Method)
	assert.Equal(t, "/_index_template/test", request.URL.Path)

	migration = opensearch.Migration{
		Method: "PUT",
		URL:    "/_plugins/_ism/policies/rollover_policy",
		Params: url.Values{
			"create": []string{
				"true",
			},
		},
		Body: map[string]any{
			"policy": map[string]any{
				"description":   "Daily rollover policy.",
				"default_state": "rollover",
				"states": []map[string]any{
					{
						"name": "rollover",
						"actions": []map[string]any{
							{
								"rollover": map[string]any{
									"min_index_age": "1d",
								},
							},
						},
						"transitions": []map[string]any{},
					},
				},
			},
			"ism_template": map[string]any{
				"index_patterns": []string{
					"import-information*",
				},
				"priority": 100,
			},
		},
	}

	request, err = migration.CreateRequest()

	assert.NoError(t, err)
	assert.Equal(t, "PUT", request.Method)
	assert.Equal(t, "/_plugins/_ism/policies/rollover_policy", request.URL.Path)
	assert.Equal(t, "create=true", request.URL.Query().Encode())
	assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

	expected := `{"ism_template":{"index_patterns":["import-information*"],"priority":100},"policy":{"default_state":"rollover","description":"Daily rollover policy.","states":[{"actions":[{"rollover":{"min_index_age":"1d"}}],"name":"rollover","transitions":[]}]}}`
	body, err := io.ReadAll(request.Body)

	assert.NoError(t, err)
	assert.Equal(t, expected, string(body))
}
