package layer4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecutorRole_String(t *testing.T) {
	tests := []struct {
		name string
		role ExecutorRole
		want string
	}{
		{
			name: "Primary role",
			role: PrimaryRole,
			want: "Primary",
		},
		{
			name: "Advisory role",
			role: AdvisoryRole,
			want: "Advisory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExecutorRole_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ExecutorRole
		wantErr bool
	}{
		{
			name:    "Primary role",
			data:    `"Primary"`,
			want:    PrimaryRole,
			wantErr: false,
		},
		{
			name:    "Advisory role",
			data:    `"Advisory"`,
			want:    AdvisoryRole,
			wantErr: false,
		},
		{
			name:    "Invalid role",
			data:    `"Invalid"`,
			want:    PrimaryRole,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var role ExecutorRole
			err := role.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, role)
			}
		})
	}
}

func TestExecutorRole_IsValid(t *testing.T) {
	tests := []struct {
		name string
		role ExecutorRole
		want bool
	}{
		{
			name: "Primary role",
			role: PrimaryRole,
			want: true,
		},
		{
			name: "Advisory role",
			role: AdvisoryRole,
			want: true,
		},
		{
			name: "Invalid role",
			role: ExecutorRole(999),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.role.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestExecutorRole_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ExecutorRole
		wantErr bool
	}{
		{
			name:    "Primary role",
			data:    "Primary",
			want:    PrimaryRole,
			wantErr: false,
		},
		{
			name:    "Advisory role",
			data:    "Advisory",
			want:    AdvisoryRole,
			wantErr: false,
		},
		{
			name:    "Invalid role",
			data:    "Invalid",
			want:    PrimaryRole,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var role ExecutorRole
			err := role.UnmarshalYAML([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, role)
			}
		})
	}
}

func TestGetEffectiveRole(t *testing.T) {
	tests := []struct {
		name      string
		mapping   ExecutorMapping
		want      ExecutorRole
		wantPanic bool
	}{
		{
			name: "Valid Primary role",
			mapping: ExecutorMapping{
				Role: PrimaryRole,
			},
			want:      PrimaryRole,
			wantPanic: false,
		},
		{
			name: "Valid Advisory role",
			mapping: ExecutorMapping{
				Role: AdvisoryRole,
			},
			want:      AdvisoryRole,
			wantPanic: false,
		},
		{
			name: "Zero value defaults to Primary",
			mapping: ExecutorMapping{
				Role: ExecutorRole(0), // Zero value
			},
			want:      PrimaryRole,
			wantPanic: false,
		},
		{
			name: "Invalid role panics",
			mapping: ExecutorMapping{
				Role: ExecutorRole(999),
			},
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() {
					GetEffectiveRole(tt.mapping)
				}, "GetEffectiveRole should panic on invalid role")
			} else {
				got := GetEffectiveRole(tt.mapping)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
