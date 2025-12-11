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
