package layer2

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
)

// loadYamlFromURL is a sub-function of loadYaml for HTTP only. It takes a URL as a sourcePath and a pointer to a Catalog object.
func loadYamlFromURL(sourcePath string, data *Catalog) error {
	resp, err := http.Get(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
	}

	err = decode(resp.Body, data)
	if err != nil {
		return fmt.Errorf("failed to decode YAML from URL: %v", err)
	}
	return nil
}

// loadYaml opens a provided path to unmarshal its data as YAML. It takes a URL or local path to a file as a sourcePath and a pointer to a Catalog object.
func loadYaml(sourcePath string, data *Catalog) error {
	if strings.HasPrefix(sourcePath, "http") {
		return loadYamlFromURL(sourcePath, data)
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

// loadJson opens a provided path to unmarshal its data as JSON. It takes a URL or local path to a file as a sourcePath and a pointer to a Catalog object.
func loadJson(sourcePath string, data *Catalog) error {
	return fmt.Errorf("loadJson not implemented [%s, %v]", sourcePath, data)
}

// LoadFiles loads data from any number of YAML files at the provided paths. JSON support is pending development.
// If run multiple times, this method will append new data to previous data.
func (c *Catalog) LoadFiles(sourcePaths []string) error {
	for _, sourcePath := range sourcePaths {
		catalog := &Catalog{}
		err := c.LoadFile(sourcePath)
		if err != nil {
			return err
		}
		c.ControlFamilies = append(c.ControlFamilies, catalog.ControlFamilies...)
		c.Capabilities = append(c.Capabilities, catalog.Capabilities...)
		c.Threats = append(c.Threats, catalog.Threats...)
	}
	return nil
}

// LoadFile loads data from a single YAML file at the provided path. JSON support is pending development.
// If run multiple times for the same data type, this method will override previous data.
func (c *Catalog) LoadFile(sourcePath string) error {
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

// decode unmarshals the provided reader into the provided Catalog object.
func decode(reader io.Reader, data *Catalog) error {
	decoder := yaml.NewDecoder(reader)
	err := decoder.Decode(data)
	if err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}
	return nil
}
