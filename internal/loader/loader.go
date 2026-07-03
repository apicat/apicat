package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3high "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func Load(dir string) (*libopenapi.DocumentModel[v3high.Document], error) {
	specPath := findSpecFile(dir)
	if specPath == "" {
		return nil, fmt.Errorf("no OpenAPI spec file found in %s (looking for openapi.yaml, openapi.yml, or openapi.json)", dir)
	}

	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read spec file %s: %w", specPath, err)
	}

	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve directory path: %w", err)
	}

	config := datamodel.NewDocumentConfiguration()
	config.BasePath = absDir
	config.AllowFileReferences = true

	doc, err := libopenapi.NewDocumentWithConfiguration(specBytes, config)
	if err != nil {
		return nil, fmt.Errorf("cannot parse OpenAPI document: %w", err)
	}

	model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("cannot build OpenAPI v3 model: %w", err)
	}
	if model == nil {
		return nil, fmt.Errorf("failed to build OpenAPI v3 model")
	}

	return model, nil
}

func findSpecFile(dir string) string {
	candidates := []string{"openapi.yaml", "openapi.yml", "openapi.json"}
	for _, name := range candidates {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}
