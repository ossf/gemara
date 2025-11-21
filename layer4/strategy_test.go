package layer4

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConflictRuleType_String(t *testing.T) {
	tests := []struct {
		name string
		rule ConflictRuleType
		want string
	}{
		{
			name: "Strict rule",
			rule: Strict,
			want: "Strict",
		},
		{
			name: "ManualOverride rule",
			rule: ManualOverride,
			want: "ManualOverride",
		},
		{
			name: "AuthoritativeConfirmation rule",
			rule: AuthoritativeConfirmation,
			want: "AuthoritativeConfirmation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.rule.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConflictRuleType_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ConflictRuleType
		wantErr bool
	}{
		{
			name:    "Strict rule",
			data:    "Strict",
			want:    Strict,
			wantErr: false,
		},
		{
			name:    "ManualOverride rule",
			data:    "ManualOverride",
			want:    ManualOverride,
			wantErr: false,
		},
		{
			name:    "AuthoritativeConfirmation rule",
			data:    "AuthoritativeConfirmation",
			want:    AuthoritativeConfirmation,
			wantErr: false,
		},
		{
			name:    "Invalid rule",
			data:    "Invalid",
			want:    Strict,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule ConflictRuleType
			err := rule.UnmarshalYAML([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, rule)
			}
		})
	}
}

func TestConflictRuleType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		want    ConflictRuleType
		wantErr bool
	}{
		{
			name:    "Strict rule",
			data:    `"Strict"`,
			want:    Strict,
			wantErr: false,
		},
		{
			name:    "ManualOverride rule",
			data:    `"ManualOverride"`,
			want:    ManualOverride,
			wantErr: false,
		},
		{
			name:    "AuthoritativeConfirmation rule",
			data:    `"AuthoritativeConfirmation"`,
			want:    AuthoritativeConfirmation,
			wantErr: false,
		},
		{
			name:    "Invalid rule",
			data:    `"Invalid"`,
			want:    Strict,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rule ConflictRuleType
			err := rule.UnmarshalJSON([]byte(tt.data))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, rule)
			}
		})
	}
}
