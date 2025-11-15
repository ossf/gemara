package layer4

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ConflictRuleType specifies the type of aggregation logic used to resolve conflicts
// when multiple executors provide results for the same assessment procedure.
// This is designed to restrict the possible conflict rule values to a set of known types.
type ConflictRuleType int

const (
	// Strict indicates that if any executor reports a failure, the overall
	// procedure result is failed, regardless of other executor results.
	Strict ConflictRuleType = iota
	// ManualOverride gives precedence to manual review executors over automated
	// executors when results conflict.
	ManualOverride
	// AdvisoryRequiresConfirmation treats Advisory executors (explicit role)
	// as requiring confirmation from Primary executors before triggering findings.
	AdvisoryRequiresConfirmation
)

var conflictRuleTypeToString = map[ConflictRuleType]string{
	Strict:                       "Strict",
	ManualOverride:               "ManualOverride",
	AdvisoryRequiresConfirmation: "AdvisoryRequiresConfirmation",
}

var stringToConflictRuleType = map[string]ConflictRuleType{
	"Strict":                       Strict,
	"ManualOverride":               ManualOverride,
	"AdvisoryRequiresConfirmation": AdvisoryRequiresConfirmation,
}

func (c *ConflictRuleType) String() string {
	return conflictRuleTypeToString[*c]
}

// MarshalYAML ensures that ConflictRuleType is serialized as a string in YAML
func (c *ConflictRuleType) MarshalYAML() (interface{}, error) {
	return c.String(), nil
}

// UnmarshalYAML ensures that ConflictRuleType can be deserialized from a YAML string
func (c *ConflictRuleType) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToConflictRuleType[s]; ok {
		*c = val
		return nil
	}
	return fmt.Errorf("invalid ConflictRuleType: %s", s)
}

// MarshalJSON ensures that ConflictRuleType is serialized as a string in JSON
func (c *ConflictRuleType) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON ensures that ConflictRuleType can be deserialized from a JSON string
func (c *ConflictRuleType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToConflictRuleType[s]; ok {
		*c = val
		return nil
	}
	return fmt.Errorf("invalid ConflictRuleType: %s", s)
}
