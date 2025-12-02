package gemara

import (
	"encoding/json"
	"fmt"

	"github.com/ossf/gemara/internal/loaders"
)

// ActorType specifies type of actor interacting in the workflow (human/software)
type ActorType int

const (
	// Software indicates the actor creates outputs without human intervention.
	Software ActorType = iota
	// Human indicates the actor creates outputs as a result of human review or judgment.
	Human
	// SoftwareAssisted indicates the actor creates outputs with software assistance but requires human oversight or judgment.
	SoftwareAssisted
)

var evaluatorTypeToString = map[ActorType]string{
	Software:        "Software",
	Human:           "Human",
	SoftwareAssisted: "Software-Assisted",
}

var stringToEvaluatorType = map[string]ActorType{
	"Software":         Software,
	"Human":            Human,
	"Software-Assisted": SoftwareAssisted,
}

func (e *ActorType) String() string {
	return evaluatorTypeToString[*e]
}

// MarshalYAML ensures that ActorType is serialized as a string in YAML
func (e *ActorType) MarshalYAML() (interface{}, error) {
	return e.String(), nil
}

// UnmarshalYAML ensures that ActorType can be deserialized from a YAML string
func (e *ActorType) UnmarshalYAML(data []byte) error {
	var s string
	if err := loaders.UnmarshalYAML(data, &s); err != nil {
		return err
	}
	if val, ok := stringToEvaluatorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ActorType: %s", s)
}

// MarshalJSON ensures that ActorType is serialized as a string in JSON
func (e *ActorType) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON ensures that ActorType can be deserialized from a JSON string
func (e *ActorType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if val, ok := stringToEvaluatorType[s]; ok {
		*e = val
		return nil
	}
	return fmt.Errorf("invalid ActorType: %s", s)
}
