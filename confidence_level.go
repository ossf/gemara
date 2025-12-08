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

const (
	// HighConfidenceThreshold is the minimum percentage (0.75 = 75%) of steps
	// that must have High confidence for the aggregated result to be High.
	HighConfidenceThreshold = 0.75
	// MediumConfidenceThreshold is the minimum percentage (0.50 = 50%) of steps
	// that must have Medium or High confidence for the aggregated result to be Medium.
	MediumConfidenceThreshold = 0.50
)

// ConfidenceAggregator tracks the distribution of confidence levels across steps
// for threshold-based aggregation.
type ConfidenceAggregator struct {
	lowCount        int
	mediumCount     int
	highCount       int
	totalSteps      int
	hasUndetermined bool
}

func NewConfidenceAggregator() *ConfidenceAggregator {
	return &ConfidenceAggregator{}
}

// Update aggregates a new confidence level using threshold-based rules and returns
// the aggregated confidence level.
func (c *ConfidenceAggregator) Update(new ConfidenceLevel) ConfidenceLevel {
	c.updateCounts(new)

	// If any step is Undetermined, result is Undetermined
	if c.hasUndetermined {
		return Undetermined
	}

	if c.totalSteps == 0 {
		return new
	}

	highPercent := float64(c.highCount) / float64(c.totalSteps)
	mediumOrHighPercent := float64(c.mediumCount+c.highCount) / float64(c.totalSteps)

	if highPercent >= HighConfidenceThreshold {
		return High
	}

	if mediumOrHighPercent >= MediumConfidenceThreshold {
		return Medium
	}

	return Low
}

func (c *ConfidenceAggregator) updateCounts(level ConfidenceLevel) {
	if level == NotSet {
		return
	}

	if level == Undetermined {
		c.hasUndetermined = true
		return
	}

	c.totalSteps++
	switch level {
	case Low:
		c.lowCount++
	case Medium:
		c.mediumCount++
	case High:
		c.highCount++
	}
}
