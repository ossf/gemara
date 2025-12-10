package gemara

import (
	"testing"
)

func TestResultString(t *testing.T) {
	tests := []struct {
		name     string
		result   Result
		expected string
	}{
		{
			result:   Passed,
			expected: "Passed",
		},
		{
			result:   Failed,
			expected: "Failed",
		},
		{
			result:   NeedsReview,
			expected: "Needs Review",
		},
		{
			result:   NotRun,
			expected: "Not Run",
		},
		{
			result:   NotApplicable,
			expected: "Not Applicable",
		},
		{
			result:   Unknown,
			expected: "Unknown",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.result.String()
			if actual != test.expected {
				t.Errorf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestUpdateAggregateResult(t *testing.T) {
	tests := []struct {
		name     string
		prev     Result
		new      Result
		expected Result
	}{
		{
			name:     "NotRun should not overwrite anything",
			prev:     Passed,
			new:      NotRun,
			expected: Passed,
		},
		{
			name:     "Failed should not be overwritten by anything",
			prev:     Failed,
			new:      Passed,
			expected: Failed,
		},
		{
			name:     "Failed should overwrite anything",
			prev:     Passed,
			new:      Failed,
			expected: Failed,
		},
		{
			name:     "Unknown should not be overwritten by NeedsReview",
			prev:     Unknown,
			new:      NeedsReview,
			expected: Unknown,
		},
		{
			name:     "Unknown should not be overwritten by Passed",
			prev:     Unknown,
			new:      Passed,
			expected: Unknown,
		},
		{
			name:     "NeedsReview should not be overwritten by Passed",
			prev:     NeedsReview,
			new:      Passed,
			expected: NeedsReview,
		},
		{
			name:     "NeedsReview should overwrite Passed",
			prev:     Passed,
			new:      NeedsReview,
			expected: NeedsReview,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := UpdateAggregateResult(test.prev, test.new)
			if actual != test.expected {
				t.Errorf("expected %s, got %s", test.expected, actual)
			}
		})
	}
}
