package opensearch

import (
	"fmt"
	"io"
	"net/http"
)

func ReadErrorFromResponse(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse response body: %w", err)
	}

	return fmt.Errorf("request failed with status code '%d': %s", resp.StatusCode, respBody)
}

func closeResponseBody(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}

	_ = resp.Body.Close()
}
