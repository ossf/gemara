package layer4

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ExecutionType specifies whether an actor's execution mode is automated or manual.
type ExecutionType int

const (
	// Automated indicates the evaluator is a tool or script that runs without human intervention.
	Automated ExecutionType = iota
	// Manual indicates the evaluator requires human review or judgment.
	Manual
)

var evaluatorTypeToString = map[ExecutionType]string{
	Automated: "Automated",
	Manual:    "Manual",
}

var stringToEvaluatorType = map[string]ExecutionType{
	"Automated": Automated,
	"Manual":    Manual,
}

func (e *ExecutionType) String() string {
	return evaluatorTypeToString[*e]
}

// MarshalYAML ensures that ExecutionType is serialized as a string in YAML
func (e *ExecutionType) MarshalYAML() (interface{}, error) {
	return e.String(), nil
}

// UnmarshalYAML ensures that ExecutionType can be deserialized from a YAML string
func (e *ExecutionType) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToEvaluatorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ExecutionType: %s", s)
}

// MarshalJSON ensures that ExecutionType is serialized as a string in JSON
func (e *ExecutionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON ensures that ExecutionType can be deserialized from a JSON string
func (e *ExecutionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToEvaluatorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ExecutionType: %s", s)
}
