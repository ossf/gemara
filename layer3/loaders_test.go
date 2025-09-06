package layer3

// This file contains table tests for the following functions:
// - loadYaml
// - loadYamlFromUri
// - loadJson (placeholder, pending implementation)
// - PolicyDocument.LoadFile

// The test data is pulled from ./test-data/

import (
	"strings"
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
		name:       "Good YAML — Policy Document",
		sourcePath: "./test-data/good-policy.yaml",
		wantErr:    false,
	},
	{
		name:       "Good YAML — Security Policy",
		sourcePath: "./test-data/good-security-policy.yml",
		wantErr:    false,
	},
}

func Test_loadYaml(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PolicyDocument{}
			if err := loadYaml(tt.sourcePath, data); (err == nil) == tt.wantErr {
				t.Errorf("loadYaml() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_LoadFile(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PolicyDocument{}
			err := p.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("PolicyDocument.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				// Validate that the policy document was loaded successfully
				assert.NotEmpty(t, p.Metadata.Id, "Policy document ID should not be empty")
				assert.NotEmpty(t, p.Metadata.Title, "Policy document title should not be empty")
				assert.NotEmpty(t, p.Metadata.Objective, "Policy document objective should not be empty")
				assert.NotEmpty(t, p.Metadata.Version, "Policy document version should not be empty")
			}
		})
	}
}

func Test_loadYamlFromUri(t *testing.T) {
	tests := []struct {
		name          string
		sourcePath    string
		wantErr       bool
		errorExpected string
	}{
		{
			name:          "URI that returns a 404",
			sourcePath:    "http://example.com/nonexistent.yaml",
			wantErr:       true,
			errorExpected: "failed to fetch Uri; response status:",
		},
		{
			name:       "Valid URI with valid data",
			sourcePath: "https://raw.githubusercontent.com/ossf/security-baseline/refs/heads/main/baseline/OSPS-AC.yaml",
			wantErr:    false,
		},
		{
			name:          "Valid URI with invalid data",
			sourcePath:    "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/template-minimum.yml",
			wantErr:       true,
			errorExpected: "failed to decode YAML from Uri:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PolicyDocument{}
			err := loadYamlFromUri(tt.sourcePath, data)
			if err != nil && tt.wantErr {
				assert.Containsf(t, err.Error(), tt.errorExpected, "expected error containing %q, got %s", tt.errorExpected, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("loadYamlFromUri() expected error matching %s, got nil.", tt.errorExpected)
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
			data := &PolicyDocument{}
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
			p := &PolicyDocument{}
			err := p.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("PolicyDocument.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_decode(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
		wantErr  bool
	}{
		{
			name: "Valid PolicyDocument YAML",
			yamlData: `
metadata:
  id: "test-policy"
  title: "Test Policy"
  objective: "Test objective"
  version: "1.0.0"
  contacts:
    author:
      name: "Test Author"
      primary: true
    responsible:
      - name: "Test Responsible"
        primary: true
    accountable:
      - name: "Test Accountable"
        primary: true
scope:
  boundaries: ["US"]
  technologies: ["cloud"]
  providers: ["aws"]
guidance-references: []
control-references: []
`,
			wantErr: false,
		},
		{
			name: "Invalid YAML structure",
			yamlData: `
metadata:
  id: "test-policy"
  title: "Test Policy"
  objective: "Test objective"
  version: "1.0.0"
  contacts:
    author:
      name: "Test Author"
      primary: true
    responsible:
      - name: "Test Responsible"
        primary: true
    accountable:
      - name: "Test Accountable"
        primary: true
scope:
  boundaries: ["US"]
  technologies: ["cloud"]
  providers: ["aws"]
guidance-references: []
control-references: []
# This should be invalid because it has an unknown field
unknown-field: "this should cause an error"
`,
			wantErr: true,
		},
		{
			name:     "Malformed YAML",
			yamlData: "this: file\nis: nonsense\n",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PolicyDocument{}
			err := decode(strings.NewReader(tt.yamlData), data)
			if (err == nil) == tt.wantErr {
				t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
