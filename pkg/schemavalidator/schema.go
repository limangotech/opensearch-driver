package schemavalidator

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

//go:embed schema.json
var schema string

type Interface interface {
	Dir(p string) error
}

type SchemaValidator struct {
}

func NewSchemaValidator() SchemaValidator {
	return SchemaValidator{}
}

func (v SchemaValidator) Dir(p string) error {
	//nolint:wrapcheck
	return filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("could not walk on provided path: %w", err)
		}

		if info.IsDir() {
			return nil
		}

		return v.File(path)
	})
}

func (v SchemaValidator) File(path string) error {
	result, err := gojsonschema.Validate(
		gojsonschema.NewStringLoader(schema),
		gojsonschema.NewReferenceLoader("file://"+path),
	)

	if err != nil {
		return fmt.Errorf("could not validate JSON schema of '%s': %w", path, err)
	}

	if !result.Valid() {
		resErrs := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			resErrs[i] = desc.String()
		}

		return fmt.Errorf("JSON schema validation failed for %s: %s", path, strings.Join(resErrs, "; "))
	}

	return nil
}
