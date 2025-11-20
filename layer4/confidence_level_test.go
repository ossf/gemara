package layer4

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

func TestConfidenceLevel_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ConfidenceLevel
		wantErr bool
	}{
		{
			name:    "Undetermined level",
			data:    "Undetermined",
			want:    Undetermined,
			wantErr: false,
		},
		{
			name:    "Low level",
			data:    "Low",
			want:    Low,
			wantErr: false,
		},
		{
			name:    "Medium level",
			data:    "Medium",
			want:    Medium,
			wantErr: false,
		},
		{
			name:    "High level",
			data:    "High",
			want:    High,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			data:    "Invalid",
			want:    Undetermined,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level ConfidenceLevel
			err := level.UnmarshalYAML([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, level)
			}
		})
	}
}

func TestConfidenceLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ConfidenceLevel
		wantErr bool
	}{
		{
			name:    "Undetermined level",
			data:    `"Undetermined"`,
			want:    Undetermined,
			wantErr: false,
		},
		{
			name:    "Low level",
			data:    `"Low"`,
			want:    Low,
			wantErr: false,
		},
		{
			name:    "Medium level",
			data:    `"Medium"`,
			want:    Medium,
			wantErr: false,
		},
		{
			name:    "High level",
			data:    `"High"`,
			want:    High,
			wantErr: false,
		},
		{
			name:    "Invalid level",
			data:    `"Invalid"`,
			want:    Undetermined,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level ConfidenceLevel
			err := level.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, level)
			}
		})
	}
}
