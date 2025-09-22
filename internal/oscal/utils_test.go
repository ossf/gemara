package oscal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeControl(t *testing.T) {
	tests := []struct {
		name        string
		inputString string
		isSubPart   bool
		wantString  string
	}{
		{
			name:        "NormalizedControl",
			inputString: "air-det-1",
			isSubPart:   false,
			wantString:  "air-det-1",
		},
		{
			name:        "CapitalizedControl",
			inputString: "AIR-DET-1",
			isSubPart:   false,
			wantString:  "air-det-1",
		},
		{
			name:        "InvalidInput/WithSubpart",
			inputString: "AU-6(9)",
			isSubPart:   false,
			wantString:  "au-6.9",
		},
		{
			name:        "Subpart",
			inputString: "AU-6(9)",
			isSubPart:   true,
			wantString:  "9",
		},
		{
			name:        "Subpart/WithMultipleDots",
			inputString: "AU-6(9).a",
			isSubPart:   true,
			wantString:  "a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotString := NormalizeControl(tt.inputString, tt.isSubPart)
			assert.Equal(t, tt.wantString, gotString)
		})
	}
}

func TestGetTimeWithFallback(t *testing.T) {
	fallbackTime := time.Date(2022, time.June, 15, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name     string
		timeStr  string
		fallback time.Time
		want     time.Time
	}{
		{
			name:     "Valid/ValidInput",
			timeStr:  "2023-01-01T12:00:00Z",
			fallback: fallbackTime,
			want:     time.Date(2023, time.January, 1, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "Valid/Fallback",
			timeStr:  "",
			fallback: fallbackTime,
			want:     fallbackTime,
		},
		{
			name:     "Invalid/Fallback",
			timeStr:  "invalid-date",
			fallback: fallbackTime,
			want:     fallbackTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetTimeWithFallback(tt.timeStr, tt.fallback)
			assert.Equal(t, tt.want, got)
		})
	}
}
