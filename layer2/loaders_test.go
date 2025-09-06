package layer2

// This file contains table tests for the following functions:
// - loadYaml
// - LoadControlFamily
// - LoadControlFamilyFiles
// - LoadControlFamiliesFile
// - loadYamlFromURL (placeholder, pending a URL to test against)
// - loadJson (placeholder, pending implementation)
// - LoadThreat (placeholder, pending implementation)
// - LoadCapability (placeholder, pending implementation)

// The test data is pulled from ./test-data.yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name       string
	sourcePath string
	wantErr    bool
}{
	{
		name:       "Bad path",
		sourcePath: "./bad-path.yaml",
		wantErr:    true,
	},
	{
		name:       "Bad YAML",
		sourcePath: "./test-data/bad.yaml",
		wantErr:    true,
	},
	{
		name:       "Good YAML — CCC",
		sourcePath: "./test-data/good-ccc.yaml",
		wantErr:    false,
	},
	{
		name:       "Good YAML — OSPS",
		sourcePath: "./test-data/good-osps.yml",
		wantErr:    false,
	},
}

func Test_loadYaml(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Catalog{}
			if err := loadYaml(tt.sourcePath, data); (err == nil) == tt.wantErr {
				t.Errorf("loadYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_LoadFile(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("Catalog.LoadControlFamily() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(c.ControlFamilies) == 0 {
				t.Errorf("Catalog.LoadControlFamily() did not load any control families")
			} else if !tt.wantErr && len(c.ControlFamilies) > 0 {
				assert.NotEmpty(t, c.ControlFamilies[0].Title, "Control family title should not be empty")
				assert.NotEmpty(t, c.ControlFamilies[0].Description, "Control family description should not be empty")
			}
		})
	}
}

func Test_LoadNestedCatalog(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadNestedCatalog(tt.sourcePath, "")
			if err == nil {
				t.Errorf("Un-nested catalogs are expected to fail")
			}
		})
	}

	nestedTests := []struct {
		name            string
		sourcePath      string
		nestedFieldName string
		wantErr         bool
	}{
		{
			name:            "Malformed URI",
			sourcePath:      "https://",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Non-conformant URI response",
			sourcePath:      "https://google.com",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Local file does not exist",
			sourcePath:      "wonky-file-name.yaml",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Empty nested catalog",
			sourcePath:      "./test-data/nested-empty.yaml",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Nested field name not provided",
			sourcePath:      "./test-data/nested-good-ccc.yaml",
			nestedFieldName: "",
			wantErr:         true,
		},
		{
			name:            "Nested field name not present in target file",
			sourcePath:      "./test-data/nested-good-ccc.yaml",
			nestedFieldName: "doesnt-exist",
			wantErr:         true,
		},
	}

	for _, tt := range nestedTests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadNestedCatalog(tt.sourcePath, tt.nestedFieldName)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error, but got none")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Did not expect error, but got '%s'", err.Error())
			} else if !tt.wantErr {
				if len(c.ControlFamilies) == 0 {
					t.Errorf("Catalog.LoadControlFamily() did not load any control families")
				} else if len(c.ControlFamilies) > 0 {
					assert.NotEmpty(t, c.ControlFamilies[0].Title, "Control family title should not be empty")
					assert.NotEmpty(t, c.ControlFamilies[0].Description, "Control family description should not be empty")
				}
			}
		})
	}
}

func Test_LoadFiles(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFiles([]string{tt.sourcePath})
			if (err == nil) == tt.wantErr {
				t.Errorf("Catalog.LoadControlFamilyFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(c.ControlFamilies) == 0 {
				t.Errorf("Catalog.LoadControlFamilyFiles() did not load any control families")
			}
		})
	}
}

func Test_loadYamlFromURL(t *testing.T) {
	tests := []struct {
		name          string
		sourcePath    string
		wantErr       bool
		errorExpected string
	}{
		{
			name:          "URL that returns a 404",
			sourcePath:    "http://example.com/nonexistent.yaml",
			wantErr:       true,
			errorExpected: "failed to fetch URL; response status:",
		},
		{
			name:       "Valid URL with valid data",
			sourcePath: "https://raw.githubusercontent.com/ossf/security-baseline/refs/heads/main/baseline/OSPS-AC.yaml",
			wantErr:    false,
		},
		{
			name:          "Valid URL with invalid data",
			sourcePath:    "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/template-minimum.yml",
			wantErr:       true,
			errorExpected: "failed to decode YAML from URL:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Catalog{}
			err := loadYamlFromURL(tt.sourcePath, data)
			if err != nil && tt.wantErr {
				assert.Containsf(t, err.Error(), tt.errorExpected, "expected error containing %q, got %s", tt.errorExpected, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("loadYamlFromURL() expected error matching %s, got nil.", tt.errorExpected)
			}
		})
	}
}

func Test_loadJson(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantErr    bool
	}{
		{
			name:       "Unsupported JSON file",
			sourcePath: "./test-data/good.json",
			wantErr:    true,
		},
		{
			name:       "Invalid JSON file",
			sourcePath: "./test-data/bad.json",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &Catalog{}
			err := loadJson(tt.sourcePath, data)
			if (err == nil) == tt.wantErr {
				t.Errorf("loadJson() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_LoadFile_UnsupportedFileType(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantErr    bool
	}{
		{
			name:       "Unsupported file type",
			sourcePath: "./test-data/unsupported.txt",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("Catalog.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
