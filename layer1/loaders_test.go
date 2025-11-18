package layer1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LoadFile(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantErr    bool
	}{
		{
			name:       "Bad path",
			sourcePath: "file://test-data/bad.yaml",
			wantErr:    true,
		},
		{
			name:       "Good YAML — AIGF",
			sourcePath: "file://test-data/good-aigf.yaml",
			wantErr:    false,
		},
		{
			name:       "Unsupported file extension",
			sourcePath: "file://test-data/unsupported.txt",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuidanceDocument{}
			err := g.LoadFile(tt.sourcePath)
			if (err == nil) == tt.wantErr {
				t.Errorf("GuidanceDocument.LoadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				assert.NotEmpty(t, g.Metadata.Id, "Guidance document ID should not be empty")
				assert.NotEmpty(t, g.Metadata.Title, "Guidance document title should not be empty")
				if len(g.Categories) == 0 {
					t.Errorf("GuidanceDocument.LoadFile() did not load any categories")
				}
			}
		})
	}
}

func Test_LoadFiles_AppendsData(t *testing.T) {
	one := &GuidanceDocument{}
	require.NoError(t, one.LoadFile("file://test-data/good-aigf.yaml"))
	require.Greater(t, len(one.Categories), 0, "expected at least one category in good-aigf.yaml")

	g := &GuidanceDocument{}
	err := g.LoadFiles([]string{
		"file://test-data/good-aigf.yaml",
		"file://test-data/good-aigf.yaml",
	})
	require.NoError(t, err)

	assert.Equal(t, one.Metadata, g.Metadata, "first document's metadata should be preserved")
	assert.Equal(t, len(one.Categories)*2, len(g.Categories), "categories should be appended across multiple files")
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
			data := &GuidanceDocument{}
			err := data.LoadFile(tt.sourcePath)
			if err != nil && tt.wantErr {
				assert.Containsf(t, err.Error(), tt.errorExpected, "expected error containing %q, got %s", tt.errorExpected, err)
			} else if err == nil && tt.wantErr {
				t.Errorf("GuidanceDocument.LoadFile() expected error matching %s, got nil.", tt.errorExpected)
			} else if err != nil && !tt.wantErr {
				t.Errorf("GuidanceDocument.LoadFile() did not expect error, but got '%s'", err.Error())
			}
		})
	}
}


