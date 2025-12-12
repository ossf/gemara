package gemara

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ResolutionStrategy specifies the type of aggregation logic used to resolve conflicts
// when multiple evaluators provide results for the same assessment procedure.
// This is designed to restrict the possible conflict rule values to a set of known types.
type ResolutionStrategy int

const (
	// MostSevere indicates that the most severe result from all evaluators is used,
	// following the severity hierarchy: Failed > Unknown > NeedsReview > Passed.
	MostSevere ResolutionStrategy = iota
	// ManualOverride gives precedence to manual review evaluators over automated
	// evaluators when results conflict.
	ManualOverride
	// AuthoritativeConfirmation treats non-authoritative evaluators
	// as requiring confirmation from authoritative evaluators before triggering findings.
	AuthoritativeConfirmation
)

var strategyToString = map[ResolutionStrategy]string{
	MostSevere:                "MostSevere",
	ManualOverride:            "ManualOverride",
	AuthoritativeConfirmation: "AuthoritativeConfirmation",
}

var stringToStrategy = map[string]ResolutionStrategy{
	"MostSevere":                MostSevere,
	"ManualOverride":            ManualOverride,
	"AuthoritativeConfirmation": AuthoritativeConfirmation,
}

func (c *ResolutionStrategy) String() string {
	return strategyToString[*c]
}

// MarshalYAML ensures that ResolutionStrategy is serialized as a string in YAML
func (c *ResolutionStrategy) MarshalYAML() (interface{}, error) {
	return c.String(), nil
}

// UnmarshalYAML ensures that ResolutionStrategy can be deserialized from a YAML string
func (c *ResolutionStrategy) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToStrategy[s]; ok {
		*c = val
		return nil
	}
	return fmt.Errorf("invalid ResolutionStrategy: %s", s)
}

// MarshalJSON ensures that ResolutionStrategy is serialized as a string in JSON
func (c *ResolutionStrategy) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON ensures that ResolutionStrategy can be deserialized from a JSON string
func (c *ResolutionStrategy) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToStrategy[s]; ok {
		*c = val
		return nil
	}
	return fmt.Errorf("invalid ResolutionStrategy: %s", s)
}
