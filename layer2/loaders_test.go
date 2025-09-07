package layer2

// This file contains table tests for the following functions:
// - loadYaml
// - LoadControlFamily
// - LoadControlFamilyFiles
// - LoadControlFamiliesFile
// - decodeYAMLFromURL (use decodeYAMLFromURL for URL-based YAML decoding)
// - loadJson (placeholder, pending implementation)
// - LoadThreat (placeholder, pending implementation)
// - LoadCapability (placeholder, pending implementation)

// The test data is pulled from ./test-data.yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	name          string
	sourcePath    string
	wantErr       bool
	errorExpected string
}{
	{
		name:       "Bad path",
		sourcePath: "file://bad-path.yaml",
		wantErr:    true,
	},
	{
		name:       "Bad YAML",
		sourcePath: "file://test-data/bad.yaml",
		wantErr:    true,
	},
	{
		name:       "Good YAML — CCC",
		sourcePath: "file://test-data/good-ccc.yaml",
		wantErr:    false,
	},
	{
		name:       "Good YAML — OSPS",
		sourcePath: "file://test-data/good-osps.yml",
		wantErr:    false,
	},
	{
		name:       "Unrecognized file extension",
		sourcePath: "file://test-data/unknown.ext",
		wantErr:    true,
	},
}

func Test_LoadFile(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("Catalog.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && len(c.ControlFamilies) == 0 {
				t.Errorf("Catalog.LoadFile() did not load any control families")
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
			sourcePath:      "file://wonky-file-name.yaml",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Empty nested catalog",
			sourcePath:      "file://test-data/nested-empty.yaml",
			nestedFieldName: "catalog",
			wantErr:         true,
		},
		{
			name:            "Nested field name present",
			sourcePath:      "file://test-data/nested-good-ccc.yaml",
			nestedFieldName: "catalog",
			wantErr:         false,
		},
		{
			name:            "Nested field name not provided",
			sourcePath:      "file://test-data/nested-good-ccc.yaml",
			nestedFieldName: "",
			wantErr:         true,
		},
		{
			name:            "Nested field name not present in target file",
			sourcePath:      "file://test-data/nested-good-ccc.yaml",
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
				assert.Equal(t, "FINOS Cloud Control Catalog", c.Metadata.Title, "Catalog title should match expected value")
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
