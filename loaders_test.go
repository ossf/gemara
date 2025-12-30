package gemara

// This file contains table-driven tests for loader functions:
// - PolicyDocument.LoadFile
// - GuidanceDocument.LoadFile and LoadFiles
// - Catalog.LoadFile, LoadFiles, and LoadNestedCatalog
//
// Test data is pulled from ./test-data/

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Policy Tests
// ============================================================================

func TestPolicy_LoadFile(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{}
			err := p.LoadFile(tt.sourcePath)

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error loading file")
				// Validate that the policy document was loaded successfully
				assert.NotEmpty(t, p.Metadata.Id, "Policy document ID should not be empty")
				assert.NotEmpty(t, p.Metadata.Version, "Policy document version should not be empty")
			}
		})
	}
}

func TestPolicyDocument_LoadFile_UnsupportedFileType(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantErr    bool
	}{
		{
			name:       "Unsupported file type",
			sourcePath: "file://test-data/unsupported.txt",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{}
			err := p.LoadFile(tt.sourcePath)

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
		})
	}
}

func TestPolicyDocument_LoadFile_URI(t *testing.T) {
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
			name:       "Valid URI with valid YAML (may not match schema)",
			sourcePath: "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/template-minimum.yml",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{}
			err := p.LoadFile(tt.sourcePath)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
				if tt.errorExpected != "" {
					assert.Contains(t, err.Error(), tt.errorExpected,
						"error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "unexpected error loading from URI")
			}
		})
	}
}

// ============================================================================
// GuidanceDocument Tests
// ============================================================================

func TestGuidanceDocument_LoadFile(t *testing.T) {
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

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error loading file")
				assert.NotEmpty(t, g.Metadata.Id, "Guidance document ID should not be empty")
				assert.NotEmpty(t, g.Title, "Guidance document title should not be empty")
				assert.NotEmpty(t, g.Families, "Guidance document should have at least one family")
				assert.NotEmpty(t, g.Guidelines, "Guidance document should have at least one guideline")
			}
		})
	}
}

func TestGuidanceDocument_LoadFiles_AppendsData(t *testing.T) {
	// Load a single file to use as baseline
	singleDoc := &GuidanceDocument{}
	require.NoError(t, singleDoc.LoadFile("file://test-data/good-aigf.yaml"))
	require.Greater(t, len(singleDoc.Families), 0,
		"expected at least one family in good-aigf.yaml")
	require.Greater(t, len(singleDoc.Guidelines), 0,
		"expected at least one guideline in good-aigf.yaml")

	// Load the same file twice to verify appending behavior
	multiDoc := &GuidanceDocument{}
	err := multiDoc.LoadFiles([]string{
		"file://test-data/good-aigf.yaml",
		"file://test-data/good-aigf.yaml",
	})
	require.NoError(t, err)

	assert.Equal(t, singleDoc.Metadata, multiDoc.Metadata,
		"first document's metadata should be preserved")
	assert.Equal(t, len(singleDoc.Families)*2, len(multiDoc.Families),
		"families should be appended across multiple files")
	assert.Equal(t, len(singleDoc.Guidelines)*2, len(multiDoc.Guidelines),
		"guidelines should be appended across multiple files")
}

func TestGuidanceDocument_LoadFile_URI(t *testing.T) {
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
			name:       "Valid URI with valid YAML (may not match schema)",
			sourcePath: "https://github.com/ossf/security-insights-spec/releases/download/v2.0.0/template-minimum.yml",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuidanceDocument{}
			err := g.LoadFile(tt.sourcePath)

			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
				if tt.errorExpected != "" {
					assert.Contains(t, err.Error(), tt.errorExpected,
						"error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "unexpected error loading from URI")
			}
		})
	}
}

// ============================================================================
// Catalog Tests
// ============================================================================

func TestCatalog_LoadFile(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFile(tt.sourcePath)

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error loading file")
				assert.NotEmpty(t, c.Families,
					"catalog should have at least one family")
				assert.NotEmpty(t, c.Controls,
					"catalog should have at least one control")
				if len(c.Families) > 0 {
					assert.NotEmpty(t, c.Families[0].Title,
						"family title should not be empty")
					assert.NotEmpty(t, c.Families[0].Description,
						"family description should not be empty")
				}
			}
		})
	}
}

func TestCatalog_LoadFile_UnsupportedFileType(t *testing.T) {
	tests := []struct {
		name       string
		sourcePath string
		wantErr    bool
	}{
		{
			name:       "Unsupported file type",
			sourcePath: "file://test-data/unsupported.txt",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFile(tt.sourcePath)

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "unexpected error")
			}
		})
	}
}

func TestCatalog_LoadFiles(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadFiles([]string{tt.sourcePath})

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error loading files")
				assert.NotEmpty(t, c.Families,
					"catalog should have at least one family")
				assert.NotEmpty(t, c.Controls,
					"catalog should have at least one control")
			}
		})
	}
}

func TestCatalog_LoadNestedCatalog(t *testing.T) {
	// Test that non-nested catalogs fail
	nonNestedTests := []struct {
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
			wantErr:    true,
		},
		{
			name:       "Good YAML — Security Policy",
			sourcePath: "file://test-data/good-security-policy.yml",
			wantErr:    true,
		},
	}

	for _, tt := range nonNestedTests {
		t.Run("Non-nested: "+tt.name, func(t *testing.T) {
			c := &Catalog{}
			err := c.LoadNestedCatalog(tt.sourcePath, "")
			assert.Error(t, err, "un-nested catalogs are expected to fail")
		})
	}

	// Test nested catalog loading
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

			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				require.NoError(t, err, "unexpected error loading nested catalog")
				assert.Equal(t, "FINOS Cloud Control Catalog", c.Title,
					"catalog title should match expected value")
				assert.NotEmpty(t, c.Families,
					"catalog should have at least one family")
				assert.NotEmpty(t, c.Controls,
					"catalog should have at least one control")
				if len(c.Families) > 0 {
					assert.NotEmpty(t, c.Families[0].Title,
						"family title should not be empty")
					assert.NotEmpty(t, c.Families[0].Description,
						"family description should not be empty")
				}
			}
		})
	}
}
