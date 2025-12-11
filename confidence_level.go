package gemara

import "encoding/json"

// ConfidenceLevel indicates the evaluator's confidence level in an assessment result.
// This is designed to restrict the possible confidence level values to a set of known levels.
type ConfidenceLevel int

const (
	// NotSet indicates the confidence level has not been set yet (initial/default state).
	NotSet ConfidenceLevel = iota
	// Undetermined indicates the confidence level could not be determined (sticky, like Unknown result).
	Undetermined
	// Low indicates the evaluator has low confidence in this result.
	Low
	// Medium indicates the evaluator has moderate confidence in this result.
	Medium
	// High indicates the evaluator has high confidence in this result.
	High
)

var confidenceLevelToString = map[ConfidenceLevel]string{
	NotSet:       "Not Set",
	Undetermined: "Undetermined",
	Low:          "Low",
	Medium:       "Medium",
	High:         "High",
}

func (c ConfidenceLevel) String() string {
	return confidenceLevelToString[c]
}

// MarshalYAML ensures that ConfidenceLevel is serialized as a string in YAML
func (c ConfidenceLevel) MarshalYAML() (interface{}, error) {
	return c.String(), nil
}

// MarshalJSON ensures that ConfidenceLevel is serialized as a string in JSON
func (c ConfidenceLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
