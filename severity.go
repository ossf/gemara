package gemara

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// Severity represents the severity level of a control.
type Severity int

const (
	SeverityUnknown Severity = iota
	SeverityInfo
	SeverityLow
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

var severityToString = map[Severity]string{
	SeverityUnknown:  "Unknown",
	SeverityInfo:     "Info",
	SeverityLow:      "Low",
	SeverityMedium:   "Medium",
	SeverityHigh:     "High",
	SeverityCritical: "Critical",
}

var stringToSeverity = map[string]Severity{
	"Unknown":  SeverityUnknown,
	"Info":     SeverityInfo,
	"Low":      SeverityLow,
	"Medium":   SeverityMedium,
	"High":     SeverityHigh,
	"Critical": SeverityCritical,
}

func (s *Severity) String() string {
	return severityToString[*s]
}

// MarshalYAML ensures that Severity is serialized as a string in YAML
func (s *Severity) MarshalYAML() (interface{}, error) {
	return s.String(), nil
}

// UnmarshalYAML ensures that Severity can be deserialized from a YAML string
func (s *Severity) UnmarshalYAML(data []byte) error {
	var str string
	if err := loaders.UnmarshalYAML(data, &str); err != nil {
		return err
	}
	if val, ok := stringToSeverity[str]; ok {
		*s = val
		return nil
	}
	return fmt.Errorf("invalid Severity: %s (valid values: None, Low, Medium, High, Critical)", str)
}

// MarshalJSON ensures that Severity is serialized as a string in JSON
func (s *Severity) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON ensures that Severity can be deserialized from a JSON string
func (s *Severity) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if val, ok := stringToSeverity[str]; ok {
		*s = val
		return nil
	}
	return fmt.Errorf("invalid Severity: %s (valid values: None, Low, Medium, High, Critical)", str)
}
