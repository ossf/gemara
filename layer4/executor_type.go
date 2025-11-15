package layer4

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ExecutorType specifies whether an executor is automated or manual.
// This is designed to restrict the possible executor type values to a set of known types.
type ExecutorType int

const (
	// Automated indicates the executor is a tool or script that runs without human intervention.
	Automated ExecutorType = iota
	// Manual indicates the executor requires human review or judgment.
	Manual
)

var executorTypeToString = map[ExecutorType]string{
	Automated: "Automated",
	Manual:    "Manual",
}

var stringToExecutorType = map[string]ExecutorType{
	"Automated": Automated,
	"Manual":    Manual,
}

func (e *ExecutorType) String() string {
	return executorTypeToString[*e]
}

// MarshalYAML ensures that ExecutorType is serialized as a string in YAML
func (e *ExecutorType) MarshalYAML() (interface{}, error) {
	return e.String(), nil
}

// UnmarshalYAML ensures that ExecutorType can be deserialized from a YAML string
func (e *ExecutorType) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToExecutorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ExecutorType: %s", s)
}

// MarshalJSON ensures that ExecutorType is serialized as a string in JSON
func (e *ExecutorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON ensures that ExecutorType can be deserialized from a JSON string
func (e *ExecutorType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToExecutorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ExecutorType: %s", s)
}
