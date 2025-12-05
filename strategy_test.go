package gemara

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolutionStrategy_String(t *testing.T) {
	tests := []struct {
		name     string
		strategy ResolutionStrategy
		expected string
	}{
		{
			name:     "MostSevere",
			strategy: MostSevere,
			expected: "MostSevere",
		},
		{
			name:     "ManualOverride",
			strategy: ManualOverride,
			expected: "ManualOverride",
		},
		{
			name:     "AuthoritativeConfirmation",
			strategy: AuthoritativeConfirmation,
			expected: "AuthoritativeConfirmation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := (&tt.strategy).String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestResolutionStrategy_MarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		strategy ResolutionStrategy
		expected string
	}{
		{
			name:     "MostSevere",
			strategy: MostSevere,
			expected: "MostSevere",
		},
		{
			name:     "ManualOverride",
			strategy: ManualOverride,
			expected: "ManualOverride",
		},
		{
			name:     "AuthoritativeConfirmation",
			strategy: AuthoritativeConfirmation,
			expected: "AuthoritativeConfirmation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := (&tt.strategy).MarshalYAML()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestResolutionStrategy_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    ResolutionStrategy
		expectError bool
	}{
		{
			name:        "Valid MostSevere",
			input:       "MostSevere",
			expected:    MostSevere,
			expectError: false,
		},
		{
			name:        "Valid ManualOverride",
			input:       "ManualOverride",
			expected:    ManualOverride,
			expectError: false,
		},
		{
			name:        "Valid AuthoritativeConfirmation",
			input:       "AuthoritativeConfirmation",
			expected:    AuthoritativeConfirmation,
			expectError: false,
		},
		{
			name:        "Invalid value",
			input:       "Invalid",
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Case sensitive - lowercase",
			input:       "mostsevere",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s ResolutionStrategy
			err := s.UnmarshalYAML([]byte(tt.input))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, s)
			}
		})
	}
}

func TestResolutionStrategy_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		strategy ResolutionStrategy
		expected string
	}{
		{
			name:     "MostSevere",
			strategy: MostSevere,
			expected: `"MostSevere"`,
		},
		{
			name:     "ManualOverride",
			strategy: ManualOverride,
			expected: `"ManualOverride"`,
		},
		{
			name:     "AuthoritativeConfirmation",
			strategy: AuthoritativeConfirmation,
			expected: `"AuthoritativeConfirmation"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := (&tt.strategy).MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestResolutionStrategy_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    ResolutionStrategy
		expectError bool
	}{
		{
			name:        "Valid MostSevere",
			input:       `"MostSevere"`,
			expected:    MostSevere,
			expectError: false,
		},
		{
			name:        "Valid ManualOverride",
			input:       `"ManualOverride"`,
			expected:    ManualOverride,
			expectError: false,
		},
		{
			name:        "Valid AuthoritativeConfirmation",
			input:       `"AuthoritativeConfirmation"`,
			expected:    AuthoritativeConfirmation,
			expectError: false,
		},
		{
			name:        "Invalid value",
			input:       `"Invalid"`,
			expectError: true,
		},
		{
			name:        "Empty string",
			input:       `""`,
			expectError: true,
		},
		{
			name:        "Case sensitive - lowercase",
			input:       `"mostsevere"`,
			expectError: true,
		},
		{
			name:        "Invalid JSON",
			input:       `not json`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s ResolutionStrategy
			err := s.UnmarshalJSON([]byte(tt.input))
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, s)
			}
		})
	}
}

func TestResolutionStrategy_RoundTrip(t *testing.T) {
	strategies := []ResolutionStrategy{
		MostSevere,
		ManualOverride,
		AuthoritativeConfirmation,
	}

	for _, original := range strategies {
		t.Run(original.String(), func(t *testing.T) {
			// Test YAML round trip
			yamlData, err := (&original).MarshalYAML()
			require.NoError(t, err)

			var yamlResult ResolutionStrategy
			err = yamlResult.UnmarshalYAML([]byte(yamlData.(string)))
			require.NoError(t, err)
			assert.Equal(t, original, yamlResult)

			// Test JSON round trip
			originalPtr := &original
			jsonData, err := json.Marshal(originalPtr)
			require.NoError(t, err)

			var jsonResult ResolutionStrategy
			err = json.Unmarshal(jsonData, &jsonResult)
			require.NoError(t, err)
			assert.Equal(t, original, jsonResult)
		})
	}
}
