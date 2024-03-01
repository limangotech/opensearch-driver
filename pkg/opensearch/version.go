package opensearch

type getDocumentResponse struct {
	Source struct {
		Version int  `json:"version"`
		Dirty   bool `json:"dirty"`
	} `json:"_source"` //nolint:tagliatelle
}
