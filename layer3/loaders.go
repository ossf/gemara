package layer3

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

func loadYamlFromUri(sourcePath string, data *PolicyDocument) error {
	resp, err := http.Get(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to fetch Uri: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch Uri; response status: %v", resp.Status)
	}

	err = decode(resp.Body, data)
	if err != nil {
		return fmt.Errorf("failed to decode YAML from Uri: %v", err)
	}
	return nil
}

func decode(reader io.Reader, data *PolicyDocument) error {
	decoder := yaml.NewDecoder(reader, yaml.DisallowUnknownField())
	err := decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}
	return nil
}

// loadYaml opens a provided path to unmarshal its data as YAML.
// sourcePath is a Uri or local path to a file.
// data is a pointer to the recieving object.
func loadYaml(sourcePath string, data *PolicyDocument) error {
	if strings.HasPrefix(sourcePath, "http") {
		return loadYamlFromUri(sourcePath, data)
	}

	file, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	err = decode(file, data)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w (%s)", err, sourcePath)
	}
	return nil
}

// loadYaml opens a provided path to unmarshal its data as JSON.
// sourcePath is a Uri or local path to a file.
// data is a pointer to the recieving object.
func loadJson(sourcePath string, data *PolicyDocument) error {
	return fmt.Errorf("loadJson not implemented [%s, %v]", sourcePath, data)
}

// LoadControlFamiliesFile loads data from a single YAML
// file at the provided path. JSON support is pending development.
// If run multiple times for the same data type, this method will override previous data.
func (c *PolicyDocument) LoadFile(sourcePath string) error {
	if strings.Contains(sourcePath, ".yaml") || strings.Contains(sourcePath, ".yml") {
		err := loadYaml(sourcePath, c)
		if err != nil {
			return err
		}
	} else if strings.Contains(sourcePath, ".json") {
		err := loadJson(sourcePath, c)
		if err != nil {
			return fmt.Errorf("error loading json: %w", err)
		}
	} else {
		return fmt.Errorf("unsupported file type")
	}
	return nil
}
