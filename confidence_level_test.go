package gemara

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfidenceLevel_String(t *testing.T) {
	tests := []struct {
		name  string
		level ConfidenceLevel
		want  string
	}{
		{
			name:  "NotSet level",
			level: NotSet,
			want:  "Not Set",
		},
		{
			name:  "Undetermined level",
			level: Undetermined,
			want:  "Undetermined",
		},
		{
			name:  "Low level",
			level: Low,
			want:  "Low",
		},
		{
			name:  "Medium level",
			level: Medium,
			want:  "Medium",
		},
		{
			name:  "High level",
			level: High,
			want:  "High",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConfidenceLevel_MarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		level   ConfidenceLevel
		want    string
		wantErr bool
	}{
		{
			name:    "NotSet level",
			level:   NotSet,
			want:    "Not Set",
			wantErr: false,
		},
		{
			name:    "Undetermined level",
			level:   Undetermined,
			want:    "Undetermined",
			wantErr: false,
		},
		{
			name:    "Low level",
			level:   Low,
			want:    "Low",
			wantErr: false,
		},
		{
			name:    "Medium level",
			level:   Medium,
			want:    "Medium",
			wantErr: false,
		},
		{
			name:    "High level",
			level:   High,
			want:    "High",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.level.MarshalYAML()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestConfidenceLevel_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		level   ConfidenceLevel
		want    string
		wantErr bool
	}{
		{
			name:    "NotSet level",
			level:   NotSet,
			want:    `"Not Set"`,
			wantErr: false,
		},
		{
			name:    "Undetermined level",
			level:   Undetermined,
			want:    `"Undetermined"`,
			wantErr: false,
		},
		{
			name:    "Low level",
			level:   Low,
			want:    `"Low"`,
			wantErr: false,
		},
		{
			name:    "Medium level",
			level:   Medium,
			want:    `"Medium"`,
			wantErr: false,
		},
		{
			name:    "High level",
			level:   High,
			want:    `"High"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.level.MarshalJSON()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestConfidenceAggregator_Update(t *testing.T) {
	tests := []struct {
		name     string
		levels   []ConfidenceLevel
		expected ConfidenceLevel
	}{
		{
			name:     "1 Low + 3 Highs → High (75% threshold)",
			levels:   []ConfidenceLevel{Low, High, High, High},
			expected: High,
		},
		{
			name:     "2 Low + 2 Highs → Medium (50% Medium+ threshold)",
			levels:   []ConfidenceLevel{Low, Low, High, High},
			expected: Medium,
		},
		{
			name:     "3 Low + 1 High → Low (<50% Medium+)",
			levels:   []ConfidenceLevel{Low, Low, Low, High},
			expected: Low,
		},
		{
			name:     "Undetermined is sticky",
			levels:   []ConfidenceLevel{High, Undetermined, High},
			expected: Undetermined,
		},
		{
			name:     "Single step",
			levels:   []ConfidenceLevel{Medium},
			expected: Medium,
		},
		{
			name:     "High + Low → Medium (50% Medium+ threshold)",
			levels:   []ConfidenceLevel{High, Low},
			expected: Medium,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			agg := NewConfidenceAggregator()
			var result ConfidenceLevel
			for _, level := range tt.levels {
				result = agg.Update(level)
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
