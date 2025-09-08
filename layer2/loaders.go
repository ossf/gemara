package layer2

import (
	"fmt"
	"path"

	"github.com/ossf/gemara/internal/loaders"
)

// LoadFiles loads data from any number of YAML or JSON files at the provided paths.
// sourcePath are expected to be file or https URIs in the form file:///path/to/file.yaml or https://example.com/file.yaml.
// If run multiple times, this method will append new data to previous data.
func (c *Catalog) LoadFiles(sourcePaths []string) error {
	for _, sourcePath := range sourcePaths {
		catalog := &Catalog{}
		err := catalog.LoadFile(sourcePath)
		if err != nil {
			return err
		}
		if c.Metadata.Id == "" {
			c.Metadata = catalog.Metadata
		}
		c.ControlFamilies = append(c.ControlFamilies, catalog.ControlFamilies...)
		c.Capabilities = append(c.Capabilities, catalog.Capabilities...)
		c.Threats = append(c.Threats, catalog.Threats...)
		c.ImportedControls = append(c.ImportedControls, catalog.ImportedControls...)
		c.ImportedCapabilities = append(c.ImportedCapabilities, catalog.ImportedCapabilities...)
		c.ImportedThreats = append(c.ImportedThreats, catalog.ImportedThreats...)
	}
	return nil
}

// LoadFile loads data from a single YAML or JSON file at the provided path.
// sourcePath is expected to be a file or https URI in the form file:///path/to/file.yaml or https://example.com/file.yaml.
// If run multiple times for the same data type, this method will override previous data.
func (c *Catalog) LoadFile(sourcePath string) error {
	ext := path.Ext(sourcePath)
	switch ext {
	case ".yaml", ".yml":
		err := loaders.LoadYAML(sourcePath, c)
		if err != nil {
			return err
		}
	case ".json":
		err := loaders.LoadJSON(sourcePath, c)
		if err != nil {
			return fmt.Errorf("error loading json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", ext)
	}
	return nil
}

// LoadNestedCatalog loads a YAML file containing a nested catalog.
// Only supports a single layer of nesting.
// Accepts file URIs with the 'file:///' prefix.
// Throws an error if the URL is not https.
// TODO: Consider validating/sanitizing inputs to reduce injection risks.
func (c *Catalog) LoadNestedCatalog(sourcePath, fieldName string) error {
	if fieldName == "" {
		return fmt.Errorf("fieldName cannot be empty")
	}
	var yamlData map[string]interface{}
	err := loaders.LoadYAML(sourcePath, &yamlData)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w (%s)", err, sourcePath)
	}
	fieldData, exists := yamlData[fieldName]
	if !exists {
		return fmt.Errorf("field '%s' not found in YAML file", fieldName)
	}
	// Marshal and unmarshal the nested field into Catalog
	fieldYamlBytes, err := loaders.MarshalYAML(fieldData)
	if err != nil {
		return fmt.Errorf("error marshaling field data to YAML: %w", err)
	}
	err = loaders.UnmarshalYAML(fieldYamlBytes, c)
	if err != nil {
		return fmt.Errorf("error decoding field '%s' into Catalog: %w", fieldName, err)
	}
	return nil
}
