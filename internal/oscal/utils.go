package oscal

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	oscalValidation "github.com/defenseunicorns/go-oscal/src/pkg/validation"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

const GemaraNamespace = "https://github.com/ossf/gemara/ns/oscal"

// NilIfEmpty returns a pointer to the slice, or nil if empty.
func NilIfEmpty[T any](slice []T) *[]T {
	if len(slice) == 0 {
		return nil
	}
	return &slice
}

// NormalizeControl alters the given control id to conform to OSCAL constraints. If the control is a
// subpart, the subpart identifier is extracted and returned.
func NormalizeControl(controlId string, subPart bool) string {
	re := regexp.MustCompile(`\((\d+)\)`)
	normalizedString := strings.ToLower(re.ReplaceAllString(controlId, ".$1"))

	if subPart {
		// This logic ensures the ids match the convention
		// <control>_<type>.<subpart>
		lastDotIndex := strings.LastIndex(normalizedString, ".")
		if lastDotIndex != -1 && lastDotIndex < len(normalizedString)-1 {
			return normalizedString[lastDotIndex+1:]
		}
	}

	return normalizedString
}

func GetTimeWithFallback(timeStr string, fallback time.Time) time.Time {
	if t := GetTime(timeStr); t != nil {
		return *t
	}
	return fallback
}

// GetTime parses a RFC3339 time string. Returns pointer to time if valid, else nil.
func GetTime(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return nil
	}
	return &t
}

func Validate(oscalModels oscal.OscalModels) error {
	validator, err := oscalValidation.NewValidatorDesiredVersion(oscalModels, oscal.Version)
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}
	if err := validator.Validate(); err != nil {
		return fmt.Errorf("model failed validation: %w", err)
	}
	return nil
}
