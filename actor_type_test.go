package gemara

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActorType_String(t *testing.T) {
	tests := []struct {
		name string
		typ  ActorType
		want string
	}{
		{
			name: "Software type",
			typ:  Software,
			want: "Software",
		},
		{
			name: "Human type",
			typ:  Human,
			want: "Human",
		},
		{
			name: "Software-Assisted type",
			typ:  SoftwareAssisted,
			want: "Software-Assisted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.typ.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestActorType_MarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		typ     ActorType
		want    string
		wantErr bool
	}{
		{
			name:    "Software type",
			typ:     Software,
			want:    "Software",
			wantErr: false,
		},
		{
			name:    "Human type",
			typ:     Human,
			want:    "Human",
			wantErr: false,
		},
		{
			name:    "Software-Assisted type",
			typ:     SoftwareAssisted,
			want:    "Software-Assisted",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.typ.MarshalYAML()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestActorType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ActorType
		wantErr bool
	}{
		{
			name:    "Software type",
			data:    "Software",
			want:    Software,
			wantErr: false,
		},
		{
			name:    "Human type",
			data:    "Human",
			want:    Human,
			wantErr: false,
		},
		{
			name:    "Software-Assisted type",
			data:    "Software-Assisted",
			want:    SoftwareAssisted,
			wantErr: false,
		},
		{
			name:    "Invalid type",
			data:    "Invalid",
			want:    Software, // Default value
			wantErr: true,
		},
		{
			name:    "Empty string",
			data:    "",
			want:    Software, // Default value
			wantErr: true,
		},
		{
			name:    "Case sensitive - lowercase",
			data:    "software",
			want:    Software, // Default value
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typ ActorType
			err := typ.UnmarshalYAML([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, Software, typ)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, typ)
			}
		})
	}
}

func TestActorType_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		typ     ActorType
		want    string
		wantErr bool
	}{
		{
			name:    "Software type",
			typ:     Software,
			want:    `"Software"`,
			wantErr: false,
		},
		{
			name:    "Human type",
			typ:     Human,
			want:    `"Human"`,
			wantErr: false,
		},
		{
			name:    "Software-Assisted type",
			typ:     SoftwareAssisted,
			want:    `"Software-Assisted"`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.typ.MarshalJSON()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}

func TestActorType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ActorType
		wantErr bool
	}{
		{
			name:    "Software type",
			data:    `"Software"`,
			want:    Software,
			wantErr: false,
		},
		{
			name:    "Human type",
			data:    `"Human"`,
			want:    Human,
			wantErr: false,
		},
		{
			name:    "Software-Assisted type",
			data:    `"Software-Assisted"`,
			want:    SoftwareAssisted,
			wantErr: false,
		},
		{
			name:    "Invalid type",
			data:    `"Invalid"`,
			want:    Software, // Default value
			wantErr: true,
		},
		{
			name:    "Empty string",
			data:    `""`,
			want:    Software, // Default value
			wantErr: true,
		},
		{
			name:    "Case sensitive - lowercase",
			data:    `"software"`,
			want:    Software, // Default value
			wantErr: true,
		},
		{
			name:    "Invalid JSON",
			data:    `not json`,
			want:    Software, // Default value
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var typ ActorType
			err := typ.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
				// When error occurs, type should remain at zero value
				assert.Equal(t, Software, typ)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, typ)
			}
		})
	}
}
