package opensearch

import (
	"crypto/tls"
	"fmt"
	"net/http"
	neturl "net/url"

	"github.com/opensearch-project/opensearch-go/v2"
)

const (
	DefaultScheme               = "https"
	DefaultPort                 = "9200"
	QueryNameInsecureSkipVerify = "insecure-skip-verify"
)

func NewTransportConfigFromURL(url string) (opensearch.Config, error) {
	parsedURL, err := neturl.Parse(url)
	if err != nil {
		return opensearch.Config{}, fmt.Errorf("failed to parse url: %w", err)
	}

	scheme := parsedURL.Scheme

	if scheme == "" {
		scheme = DefaultScheme
	}

	baseURL := scheme + "://" + parsedURL.Hostname()
	port := parsedURL.Port()

	if port == "" {
		port = DefaultPort
	}

	baseURL += ":" + port

	config := opensearch.Config{
		Addresses: []string{baseURL},
		Username:  parsedURL.User.Username(),
	}

	passwd, isSet := parsedURL.User.Password()
	if isSet {
		config.Password = passwd
	}

	if parsedURL.Query().Get(QueryNameInsecureSkipVerify) == "true" {
		config.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		}
	}

	return config, nil
}
