package layer1

import (
	"fmt"
	"path"

	"github.com/ossf/gemara/internal/loaders"
)

// LoadFiles loads data from any number of YAML or JSON files at the provided paths.
// sourcePath are expected to be file or https URIs in the form file:///path/to/file.yaml or https://example.com/file.yaml.
// If run multiple times, this method will append new data to previous data.
func (g *GuidanceDocument) LoadFiles(sourcePaths []string) error {
	for _, sourcePath := range sourcePaths {
		doc := &GuidanceDocument{}
		if err := doc.LoadFile(sourcePath); err != nil {
			return err
		}
		if g.Metadata.Id == "" {
			g.Metadata = doc.Metadata
		}
		g.Categories = append(g.Categories, doc.Categories...)
		g.ImportedGuidelines = append(g.ImportedGuidelines, doc.ImportedGuidelines...)
		g.ImportedPrinciples = append(g.ImportedPrinciples, doc.ImportedPrinciples...)
	}
	return nil
}

// LoadFile loads data from a YAML or JSON file at the provided path into the GuidanceDocument.
// sourcePath is expected to be a file or https URI in the form file:///path/to/file.yaml or https://example.com/file.yaml.
// If run multiple times for the same data type, this method will override previous data.
func (g *GuidanceDocument) LoadFile(sourcePath string) error {
	ext := path.Ext(sourcePath)
	switch ext {
	case ".yaml", ".yml":
		err := loaders.LoadYAML(sourcePath, g)
		if err != nil {
			return err
		}
	case ".json":
		err := loaders.LoadJSON(sourcePath, g)
		if err != nil {
			return fmt.Errorf("error loading json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	return nil
}


