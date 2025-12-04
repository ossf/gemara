package gemara

import "encoding/json"

// Result is an enum representing the result of a control evaluation
// This is designed to restrict the possible result values to a set of known states
type Result int

const (
	NotRun Result = iota
	Passed
	Failed
	NeedsReview
	NotApplicable
	Unknown
)

var toString = map[Result]string{
	NotRun:        "Not Run",
	Passed:        "Passed",
	Failed:        "Failed",
	NeedsReview:   "Needs Review",
	NotApplicable: "Not Applicable",
	Unknown:       "Unknown",
}

func (r Result) String() string {
	return toString[r]
}

// MarshalYAML ensures that Result is serialized as a string in YAML
func (r Result) MarshalYAML() (interface{}, error) {
	return r.String(), nil
}

// MarshalJSON ensures that Result is serialized as a string in JSON
func (r Result) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UpdateAggregateResult compares the current result with the new result and returns the most severe of the two.
func UpdateAggregateResult(previous Result, new Result) Result {
	if new == NotRun {
		// Not Run should not overwrite anything
		// Failed should not be overwritten by anything
		// Failed should overwrite anything
		return previous
	}

	if previous == Failed || new == Failed {
		// Failed should not be overwritten by anything
		// Failed should overwrite anything
		return Failed
	}

	if previous == Unknown || new == Unknown {
		// If the current or past result is Unknown, it should not be overwritten by NeedsReview or Passed.
		return Unknown
	}

	if previous == NeedsReview || new == NeedsReview {
		// NeedsReview should not be overwritten by Passed
		// NeedsReview should overwrite Passed
		return NeedsReview
	}
	return Passed
}
