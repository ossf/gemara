package layer4

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ExecutorRole determines how an executor participates in conflict resolution
// when using AdvisoryRequiresConfirmation strategy.
type ExecutorRole int

const (
	// PrimaryRole indicates the executor can trigger findings independently.
	PrimaryRole ExecutorRole = iota
	// AdvisoryRole indicates the executor requires confirmation from Primary executors
	// to trigger findings.
	AdvisoryRole
)

var roleToString = map[ExecutorRole]string{
	PrimaryRole:  "Primary",
	AdvisoryRole: "Advisory",
}

var stringToExecutorRole = map[string]ExecutorRole{
	"Primary":  PrimaryRole,
	"Advisory": AdvisoryRole,
}

// String returns the string representation of the executor role.
func (r *ExecutorRole) String() string {
	return roleToString[*r]
}

// MarshalYAML ensures that ExecutorRole is serialized as a string in YAML
func (r *ExecutorRole) MarshalYAML() (interface{}, error) {
	return r.String(), nil
}

// UnmarshalYAML ensures that ExecutorRole can be deserialized from a YAML string
func (r *ExecutorRole) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToExecutorRole[s]; ok {
		*r = val
		return nil
	}
	return fmt.Errorf("invalid ExecutorRole: %s", s)
}

// MarshalJSON ensures that ExecutorRole is serialized as a string in JSON
func (r *ExecutorRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON ensures that ExecutorRole can be deserialized from a JSON string
func (r *ExecutorRole) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToExecutorRole[s]; ok {
		*r = val
		return nil
	}
	return fmt.Errorf("invalid ExecutorRole: %s", s)
}

// IsValid checks if the executor role is valid.
func (r *ExecutorRole) IsValid() bool {
	return *r == PrimaryRole || *r == AdvisoryRole
}

// GetEffectiveRole returns the effective executor role for an ExecutorMapping.
// If the role is not explicitly set (zero value), it defaults to PrimaryRole.
// If an invalid role value is provided, it panics to prevent silent failures.
func GetEffectiveRole(mapping ExecutorMapping) ExecutorRole {
	if mapping.Role.IsValid() {
		return mapping.Role
	}

	// Invalid role value - this is a programming error and should not be silently ignored
	panic(fmt.Sprintf("invalid ExecutorRole value: %d (valid values are %d=Primary, %d=Advisory)",
		mapping.Role, PrimaryRole, AdvisoryRole))
}
