package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	ErrInvalidContent = errors.New("invalid migration content")
)

type Migration struct {
	Method  string         `json:"method"`
	URL     string         `json:"url"`
	Params  url.Values     `json:"params,omitempty"`
	Body    map[string]any `json:"body,omitempty"`
	Headers http.Header    `json:"headers,omitempty"`
}

func NewMigrationFromRawContent(content []byte) (Migration, error) {
	mgr := Migration{}
	if err := json.Unmarshal(content, &mgr); err != nil {
		return mgr, fmt.Errorf("JSON syntax error: %w", err)
	}

	if err := mgr.Validate(); err != nil {
		return mgr, fmt.Errorf("migration content validation failure: %w", err)
	}

	return mgr, nil
}

func (m Migration) Validate() error {
	switch {
	case m.URL == "":
	case m.Method == "":
		return ErrInvalidContent
	}

	return nil
}

func (m Migration) CreateRequest() (*http.Request, error) {
	body, err := json.Marshal(m.Body)
	if err != nil {
		return nil, fmt.Errorf("could not convert body to JSON: %w", err)
	}

	req, err := http.NewRequestWithContext(context.Background(), strings.ToUpper(m.Method), m.URL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request: %w", err)
	}

	req.Header = m.Headers
	if req.Header == nil {
		req.Header = http.Header{}
	}

	req.Header.Add("Content-Type", "application/json")
	setURLParams(req, m.Params)

	return req, nil
}

func setURLParams(req *http.Request, params url.Values) {
	if len(params) == 0 {
		return
	}

	query := req.URL.Query()

	for k, v := range params {
		for _, j := range v {
			query.Set(k, j)
		}
	}

	req.URL.RawQuery = query.Encode()
}
