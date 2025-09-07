package loaders

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"log"

	"github.com/goccy/go-yaml"
)

// LoadYAML loads YAML from a file or URL into the provided target struct.
func LoadYAML(sourcePath string, target interface{}) error {
	parsedURL, err := url.Parse(sourcePath)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "https" {
		return decodeYAMLFromURL(parsedURL.String(), target)
	}
	if parsedURL.Scheme == "file" {
		return decodeYAMLFromFile(strings.TrimPrefix(parsedURL.String(), "file://"), target)
	}
	return fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
}

// LoadJSON loads JSON from a file or URL into the provided target struct.
func LoadJSON(sourcePath string, target interface{}) error {
	parsedURL, err := url.Parse(sourcePath)
	if err != nil {
		return err
	}
	if parsedURL.Scheme == "https" {
		return decodeJSONFromURL(parsedURL.String(), target)
	}
	if parsedURL.Scheme == "file" {
		return decodeJSONFromFile(strings.TrimPrefix(parsedURL.String(), "file://"), target)
	}
	return fmt.Errorf("unsupported scheme: %s", parsedURL.Scheme)
}

func decodeYAMLFromURL(urlStr string, target interface{}) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
	}
	return decodeYAMLFromReader(resp.Body, target)
}

func decodeYAMLFromFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	return decodeYAMLFromReader(file, target)
}

func decodeYAMLFromReader(reader io.Reader, target interface{}) error {
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding YAML: %w", err)
	}
	return nil
}

func decodeJSONFromURL(urlStr string, target interface{}) error {
	resp, err := http.Get(urlStr)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch URL; response status: %v", resp.Status)
	}
	return decodeJSONFromReader(resp.Body, target)
}

func decodeJSONFromFile(filePath string, target interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	return decodeJSONFromReader(file, target)
}

func decodeJSONFromReader(reader io.Reader, target interface{}) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("error decoding JSON: %w", err)
	}
	return nil
}

// MarshalYAML marshals an object to YAML bytes.
func MarshalYAML(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

// UnmarshalYAML unmarshals YAML bytes into the provided target.
func UnmarshalYAML(data []byte, target interface{}) error {
	return yaml.Unmarshal(data, target)
}
