package layer3

// This file contains table tests for the following functions:
// - PolicyDocument.LoadFile

// The test data is pulled from ./test-data/

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
		sourcePath: "file://bad-path.yaml",
		wantErr:    true,
	},
	{
		name:       "Bad YAML",
		sourcePath: "file://test-data/bad.yaml",
		wantErr:    true,
	},
	{
		name:       "Good YAML — Policy Document",
		sourcePath: "file://test-data/good-policy.yaml",
		wantErr:    false,
	},
	{
		name:       "Good YAML — Security Policy",
		sourcePath: "file://test-data/good-security-policy.yml",
		wantErr:    false,
	},
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

func Test_LoadFile_Uri(t *testing.T) {
	tests := []struct {
		name          string
		sourcePath    string
		wantErr       bool
		errorExpected string
	}{
		{
			name:          "URI that returns a 404",
			sourcePath:    "https://example.com/nonexistent.yaml",
			wantErr:       true,
			errorExpected: "failed to fetch URL; response status: 404 Not Found",
		},
		{
			name:       "Valid URI with valid data",
			sourcePath: "https://raw.githubusercontent.com/ossf/security-baseline/refs/heads/main/baseline/OSPS-AC.yaml",
			wantErr:    false,
		},
		{
			name:       "Valid URI with invalid data",
			sourcePath: "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/template-minimum.yml",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := &PolicyDocument{}
			err := data.LoadFile(tt.sourcePath)
			if err != nil && tt.wantErr {
				assert.Containsf(t, err.Error(), tt.errorExpected, "expected error containing %q, got %s", tt.errorExpected, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("loadYamlFromUri() expected error matching %s, got nil.", tt.errorExpected)
			}
		})
	}
}
