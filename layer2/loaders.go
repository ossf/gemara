package layer2

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-yaml"
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
		c.ControlFamilies = append(c.ControlFamilies, catalog.ControlFamilies...)
		c.Capabilities = append(c.Capabilities, catalog.Capabilities...)
		c.Threats = append(c.Threats, catalog.Threats...)
	}
	return nil
}

// LoadFile loads data from a single YAML or JSON file at the provided path.
// sourcePath is expected to be a file or https URI in the form file:///path/to/file.yaml or https://example.com/file.yaml.
// If run multiple times for the same data type, this method will override previous data.
func (c *Catalog) LoadFile(sourcePath string) error {
	parsedURL, err := parsePath(sourcePath)
	if err != nil {
		return err
	}
	switch path.Ext(parsedURL.Path) {
	case ".yaml", ".yml":
		err := loadYaml(parsedURL, c)
		if err != nil {
			return err
		}
	case ".json":
		err := loadJson(parsedURL, c)
		if err != nil {
			return fmt.Errorf("error loading json: %w", err)
		}
	default:
		return fmt.Errorf("unsupported file extension: %s", path.Ext(parsedURL.Path))
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

	parsedURL, err := parsePath(sourcePath)
	if err != nil {
		return err
	}

	switch parsedURL.Scheme {
	case "https":
		err := decodeYAMLFromURL(parsedURL, &yamlData)
		if err != nil {
			return fmt.Errorf("error decoding YAML: %w (%s)", err, sourcePath)
		}
	case "file":
		err := decodeYAMLFromFile(parsedURL, &yamlData)
		if err != nil {
			return fmt.Errorf("error decoding YAML: %w (%s)", err, sourcePath)
		}
	default:
		return fmt.Errorf("unsupported sourcePath scheme in %s: %s", parsedURL.Path, parsedURL.Scheme)
	}

	// TODO: Validate fieldName and the nested YAML content for injection risks.
	fieldData, exists := yamlData[fieldName]
	if !exists {
		return fmt.Errorf("field '%s' not found in YAML file", fieldName)
	}
	fieldYamlBytes, err := yaml.Marshal(fieldData)
	if err != nil {
		return fmt.Errorf("error marshaling field data to YAML: %w", err)
	}
	err = yaml.Unmarshal(fieldYamlBytes, c)
	if err != nil {
		return fmt.Errorf("error decoding field '%s' into Catalog: %w", fieldName, err)
	}

	return nil
}

// decodeYAMLFromReader decodes YAML from an io.Reader into the provided target.
func decodeYAMLFromReader(reader io.Reader, target interface{}) error {
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}
	return nil
}

// decodeYAMLFromURL fetches a URL and decodes YAML into the provided target.
func decodeYAMLFromURL(url *url.URL, target interface{}) error {
	resp, err := http.Get(url.String())
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close() //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
	}
	return decodeYAMLFromReader(resp.Body, target)
}

// decodeYAMLFromFile opens a file and decodes YAML into the provided target.
func decodeYAMLFromFile(url *url.URL, target interface{}) error {
	file, err := os.Open(strings.TrimPrefix(url.String(), "file://"))
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close() //nolint:errcheck
	return decodeYAMLFromReader(file, target)
}

// loadYaml opens a provided path to unmarshal its data as YAML. It takes a URL or local path to a file as a sourcePath and a pointer to a Catalog object.
func loadYaml(url *url.URL, data *Catalog) error {
	switch url.Scheme {
	case "https":
		return decodeYAMLFromURL(url, data)
	case "file":
		return decodeYAMLFromFile(url, data)
	}
	return fmt.Errorf("loadYaml not implemented [%s, %v]", url.String(), data)
}

// loadJson opens a provided path to unmarshal its data as JSON. It takes a URL or local path to a file as a sourcePath and a pointer to a Catalog object.
func loadJson(url *url.URL, data *Catalog) error {
	switch url.Scheme {
	case "https":
		resp, err := http.Get(url.String())
		if err != nil {
			return fmt.Errorf("failed to fetch URL: %v", err)
		}
		defer resp.Body.Close() //nolint:errcheck
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
		}
		return decodeJSONFromReader(resp.Body, data)
	case "file":
		file, err := os.Open(strings.TrimPrefix(url.String(), "file://"))
		if err != nil {
			return fmt.Errorf("error opening file: %w", err)
		}
		defer file.Close() //nolint:errcheck
		return decodeJSONFromReader(file, data)
	}
	return fmt.Errorf("loadJson not implemented [%s, %v]", url.String(), data)
}

// decodeJSONFromReader decodes JSON from an io.Reader into the provided target.
func decodeJSONFromReader(reader io.Reader, target interface{}) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}

func parsePath(sourcePath string) (*url.URL, error) {
	url, error := url.Parse(sourcePath)
	if error != nil {
		return url, fmt.Errorf("failed to parse sourcePath: %w", error)
	}
	return validateURL(url)
}

var VALID_EXTENSIONS = []string{".yaml", ".yml", ".json"}

func hasValidExtension(ext string) bool {
	for _, validExt := range VALID_EXTENSIONS {
		if ext == validExt {
			return true
		}
	}
	return false
}

func validateURL(parsedURL *url.URL) (*url.URL, error) {
	switch parsedURL.Scheme {
	case "https", "file":
		if hasValidExtension(path.Ext(parsedURL.Path)) {
			return parsedURL, nil
		} else {
			return nil, fmt.Errorf("unsupported file type in URL: %s", parsedURL.Path)
		}
	default:
		return nil, fmt.Errorf("unsupported sourcePath scheme in %s: %s", parsedURL.Path, parsedURL.Scheme)
	}
}
