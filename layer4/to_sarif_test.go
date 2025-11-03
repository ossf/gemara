package layer4

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToSARIF(t *testing.T) {
	var sarif *SarifReport
	// Create test steps for LogicalLocation testing
	step1 := func(interface{}) (Result, string) { return Failed, "" }
	step2 := func(interface{}) (Result, string) { return NeedsReview, "" }
	step3 := func(interface{}) (Result, string) { return Passed, "" }

	ce := &ControlEvaluation{
		Name: "Example Control",
		Control: Mapping{
			EntryId: "CTRL-1",
		},
		Result: Passed,
		AssessmentLogs: []*AssessmentLog{
			{
				Requirement: Mapping{
					EntryId: "REQ-1",
				},
				Description:   "should do a thing",
				Result:        Failed,
				Message:       "thing was not done",
				Steps:         []AssessmentStep{step1},
				StepsExecuted: 1,
			},
			{
				Requirement: Mapping{
					EntryId: "REQ-2",
				},
				Description:   "should maybe do a thing",
				Result:        NeedsReview,
				Steps:         []AssessmentStep{step2},
				StepsExecuted: 1,
			},
			{
				Requirement: Mapping{
					EntryId: "REQ-3",
				},
				Description:   "should do another thing",
				Result:        Passed,
				Steps:         []AssessmentStep{step3},
				StepsExecuted: 1,
			},
		},
	}
	informationURI := "https://github.com/ossf/gemara"
	version := "1.0.0"

	evaluationLog := EvaluationLog{
		Evaluations: []*ControlEvaluation{ce},
		Metadata: Metadata{
			Author: Author{
				Name:    "gemara",
				Uri:     informationURI,
				Version: version,
			},
		},
	}
	// Test without parameter - should use Metadata.Author.Uri (backward compatible)
	sarifBytes, err := evaluationLog.ToSARIF("")
	require.NoError(t, err)
	sarif = &SarifReport{}
	err = json.Unmarshal(sarifBytes, sarif)
	require.NoError(t, err)
	require.NotNil(t, sarif)
	require.Len(t, sarif.Runs, 1)
	run := sarif.Runs[0]

	// rules should be unique for each requirement
	require.NotNil(t, run.Tool.Driver.Rules)
	require.Len(t, run.Tool.Driver.Rules, 3)

	// results should be present with appropriate levels
	require.Len(t, run.Results, 3)
	// map of ruleId to level
	levels := map[string]string{}
	for _, r := range run.Results {
		levels[r.RuleID] = r.Level
	}
	require.Equal(t, "error", levels["CTRL-1/REQ-1"])   // Failed
	require.Equal(t, "warning", levels["CTRL-1/REQ-2"]) // NeedsReview
	require.Equal(t, "note", levels["CTRL-1/REQ-3"])    // Passed

	// Check tool version information
	require.Equal(t, "gemara", run.Tool.Driver.Name)
	require.Equal(t, informationURI, run.Tool.Driver.InformationURI)
	require.Equal(t, version, run.Tool.Driver.Version)

	// Check that PhysicalLocation is included in all results (uses Metadata.Author.Uri as fallback)
	for _, result := range run.Results {
		require.NotEmpty(t, result.Locations, "Each result should have at least one location")
		for _, location := range result.Locations {
			require.NotNil(t, location.PhysicalLocation, "PhysicalLocation should be present")
			require.Equal(t, informationURI, location.PhysicalLocation.ArtifactLocation.URI, "ArtifactLocation URI should match Metadata.Author.Uri")
			require.Nil(t, location.PhysicalLocation.Region, "Region should be nil")
			// LogicalLocations should be present with AssessmentStep function name
			require.NotEmpty(t, location.LogicalLocations, "LogicalLocations should still be present")
		}
	}

	// ensure JSON marshals cleanly
	_, err = json.Marshal(sarif)
	require.NoError(t, err)
}

func Test_ToSARIF_NoPhysicalLocationWhenURIMissing(t *testing.T) {
	// Test that PhysicalLocation is nil when Metadata.Author.Uri is empty
	ce := &ControlEvaluation{
		Name: "Example Control",
		Control: Mapping{
			EntryId: "CTRL-1",
		},
		Result: Passed,
		AssessmentLogs: []*AssessmentLog{
			{
				Requirement: Mapping{
					EntryId: "REQ-1",
				},
				Description:   "should do a thing",
				Result:        Failed,
				Message:       "thing was not done",
				Steps:         []AssessmentStep{func(interface{}) (Result, string) { return Failed, "" }},
				StepsExecuted: 1,
			},
		},
	}

	evaluationLog := EvaluationLog{
		Evaluations: []*ControlEvaluation{ce},
		Metadata: Metadata{
			Author: Author{
				Name:    "gemara",
				Uri:     "", // Empty URI
				Version: "1.0.0",
			},
		},
	}
	// Test without parameter and empty URI - PhysicalLocation should be nil
	sarifBytes, err := evaluationLog.ToSARIF("")
	require.NoError(t, err)
	sarif := &SarifReport{}
	err = json.Unmarshal(sarifBytes, sarif)
	require.NoError(t, err)
	require.NotNil(t, sarif)
	require.Len(t, sarif.Runs, 1)
	run := sarif.Runs[0]

	// When URI is empty, PhysicalLocation should be nil
	require.Len(t, run.Results, 1)
	result := run.Results[0]
	require.NotEmpty(t, result.Locations)
	location := result.Locations[0]
	require.Nil(t, location.PhysicalLocation, "PhysicalLocation should be nil when Metadata.Author.Uri is empty")
	// LogicalLocations should still be present
	require.NotEmpty(t, location.LogicalLocations, "LogicalLocations should still be present")
}

func Test_ToSARIF_WithArtifactURIParameter(t *testing.T) {
	// Test that provided artifactURI parameter is used instead of Metadata.Author.Uri
	testStep := func(interface{}) (Result, string) {
		return Failed, ""
	}
	ce := &ControlEvaluation{
		Name: "Example Control",
		Control: Mapping{
			EntryId: "CTRL-1",
		},
		Result: Passed,
		AssessmentLogs: []*AssessmentLog{
			{
				Requirement: Mapping{
					EntryId: "REQ-1",
				},
				Description:   "Test requirement",
				Result:        Failed,
				Message:       "Test message",
				Steps:         []AssessmentStep{testStep},
				StepsExecuted: 1,
			},
		},
	}

	customURI := "README.md"
	evaluationLog := EvaluationLog{
		Evaluations: []*ControlEvaluation{ce},
		Metadata: Metadata{
			Author: Author{
				Name:    "gemara",
				Uri:     "https://github.com/test/repo", // This should NOT be used
				Version: "1.0.0",
			},
		},
	}

	// Test with custom artifactURI parameter
	sarifBytes, err := evaluationLog.ToSARIF(customURI)
	require.NoError(t, err)
	sarif := &SarifReport{}
	err = json.Unmarshal(sarifBytes, sarif)
	require.NoError(t, err)
	require.NotNil(t, sarif)
	require.Len(t, sarif.Runs, 1)
	run := sarif.Runs[0]

	require.Len(t, run.Results, 1)
	result := run.Results[0]
	require.NotEmpty(t, result.Locations)
	location := result.Locations[0]

	// Verify custom URI is used (not Metadata.Author.Uri)
	require.NotNil(t, location.PhysicalLocation)
	require.Equal(t, customURI, location.PhysicalLocation.ArtifactLocation.URI, "Should use provided artifactURI parameter")
	require.NotEqual(t, "https://github.com/test/repo", location.PhysicalLocation.ArtifactLocation.URI, "Should not use Metadata.Author.Uri when parameter provided")
}
