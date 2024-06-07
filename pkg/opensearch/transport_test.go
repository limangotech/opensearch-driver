package opensearch_test

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/stretchr/testify/assert"

	opensearchdriver "github.com/limangotech/opensearch-driver/pkg/opensearch"
)

func TestItParsesURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		url      string
		expected opensearch.Config
	}{
		{
			// https connection
			url: "https://admin:password@opensearch:1234",
			expected: opensearch.Config{
				Addresses: []string{"https://opensearch:1234"},
				Username:  "admin",
				Password:  "password",
			},
		},
		{
			// http connection
			url: "http://admin:password@opensearch:1234",
			expected: opensearch.Config{
				Addresses: []string{"http://opensearch:1234"},
				Username:  "admin",
				Password:  "password",
			},
		},
		{
			// connection with no password
			url: "http://admin@opensearch:1234",
			expected: opensearch.Config{
				Addresses: []string{"http://opensearch:1234"},
				Username:  "admin",
			},
		},
		{
			// connection with no auth
			url:      "https://opensearch:1234",
			expected: opensearch.Config{Addresses: []string{"https://opensearch:1234"}},
		},
		{
			// connection with no port
			url: "https://admin:password@opensearch",
			expected: opensearch.Config{
				Addresses: []string{"https://opensearch:9200"},
				Username:  "admin",
				Password:  "password",
			},
		},
		{
			// connection with no TLS check
			url: "https://admin:password@opensearch?insecure-skip-verify=true",
			expected: opensearch.Config{
				Addresses: []string{"https://opensearch:9200"},
				Username:  "admin",
				Password:  "password",
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			},
		},
	}

	for i, testCase := range testCases {
		actual, err := opensearchdriver.NewTransportConfigFromURL(testCase.url)

		assert.NoError(t, err, fmt.Sprintf("Case %d", i))
		assert.Equal(t, testCase.expected, actual, fmt.Sprintf("Case %d", i))
	}
}
