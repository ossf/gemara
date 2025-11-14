package layer4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutorType_String(t *testing.T) {
	tests := []struct {
		name string
		typ  ExecutorType
		want string
	}{
		{
			name: "Automated type",
			typ:  Automated,
			want: "Automated",
		},
		{
			name: "Manual type",
			typ:  Manual,
			want: "Manual",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typ.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExecutorType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ExecutorType
		wantErr bool
	}{
		{
			name:    "Automated type",
			data:    "Automated",
			want:    Automated,
			wantErr: false,
		},
		{
			name:    "Manual type",
			data:    "Manual",
			want:    Manual,
			wantErr: false,
		},
		{
			name:    "Invalid type",
			data:    "Invalid",
			want:    Automated,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typ ExecutorType
			err := typ.UnmarshalYAML([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, typ)
			}
		})
	}
}

func TestExecutorType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ExecutorType
		wantErr bool
	}{
		{
			name:    "Automated type",
			data:    `"Automated"`,
			want:    Automated,
			wantErr: false,
		},
		{
			name:    "Manual type",
			data:    `"Manual"`,
			want:    Manual,
			wantErr: false,
		},
		{
			name:    "Invalid type",
			data:    `"Invalid"`,
			want:    Automated,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typ ExecutorType
			err := typ.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, typ)
			}
		})
	}
}
