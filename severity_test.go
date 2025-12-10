package gemara

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeverityString(t *testing.T) {
	tests := []struct {
		name     string
		severity Severity
		expected string
	}{
		{
			name:     "Unknown",
			severity: SeverityUnknown,
			expected: "Unknown",
		},
		{
			name:     "Info",
			severity: SeverityInfo,
			expected: "Info",
		},
		{
			name:     "Low",
			severity: SeverityLow,
			expected: "Low",
		},
		{
			name:     "Medium",
			severity: SeverityMedium,
			expected: "Medium",
		},
		{
			name:     "High",
			severity: SeverityHigh,
			expected: "High",
		},
		{
			name:     "Critical",
			severity: SeverityCritical,
			expected: "Critical",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.severity.String()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestSeverityMarshalYAML(t *testing.T) {
	severity := SeverityHigh
	result, err := severity.MarshalYAML()
	require.NoError(t, err)
	assert.Equal(t, "High", result)
}

func TestSeverityUnmarshalYAML(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Severity
		expectError bool
	}{
		{
			name:        "Valid Unknown",
			input:       "Unknown",
			expected:    SeverityUnknown,
			expectError: false,
		},
		{
			name:        "Valid Info",
			input:       "Info",
			expected:    SeverityInfo,
			expectError: false,
		},
		{
			name:        "Valid Low",
			input:       "Low",
			expected:    SeverityLow,
			expectError: false,
		},
		{
			name:        "Valid Medium",
			input:       "Medium",
			expected:    SeverityMedium,
			expectError: false,
		},
		{
			name:        "Valid High",
			input:       "High",
			expected:    SeverityHigh,
			expectError: false,
		},
		{
			name:        "Valid Critical",
			input:       "Critical",
			expected:    SeverityCritical,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Severity
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

func TestSeverityMarshalJSON(t *testing.T) {
	severity := SeverityCritical
	data, err := severity.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `"Critical"`, string(data))
}

func TestSeverityUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    Severity
		expectError bool
	}{
		{
			name:        "Valid High",
			input:       `"High"`,
			expected:    SeverityHigh,
			expectError: false,
		},
		{
			name:        "Valid Medium",
			input:       `"Medium"`,
			expected:    SeverityMedium,
			expectError: false,
		},
		{
			name:        "Invalid value",
			input:       `"Invalid"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Severity
			err := json.Unmarshal([]byte(tt.input), &s)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, s)
			}
		})
	}
}

func TestSeverityRoundTrip(t *testing.T) {
	severities := []Severity{
		SeverityUnknown,
		SeverityInfo,
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}

	for _, original := range severities {
		t.Run(original.String(), func(t *testing.T) {
			// Test YAML round trip
			yamlData, err := original.MarshalYAML()
			require.NoError(t, err)

			var yamlResult Severity
			err = yamlResult.UnmarshalYAML([]byte(yamlData.(string)))
			require.NoError(t, err)
			assert.Equal(t, original, yamlResult)

			// Test JSON round trip (using pointer to ensure MarshalJSON is called)
			originalPtr := &original
			jsonData, err := json.Marshal(originalPtr)
			require.NoError(t, err)

			var jsonResult Severity
			err = json.Unmarshal(jsonData, &jsonResult)
			require.NoError(t, err)
			assert.Equal(t, original, jsonResult)
		})
	}
}
