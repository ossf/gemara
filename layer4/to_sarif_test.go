package layer4

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToSARIF(t *testing.T) {
	var sarif *SarifReport
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
				Description: "should do a thing",
				Result:      Failed,
				Message:     "thing was not done",
			},
			{
				Requirement: Mapping{
					EntryId: "REQ-2",
				},
				Description: "should maybe do a thing",
				Result:      NeedsReview,
			},
			{
				Requirement: Mapping{
					EntryId: "REQ-3",
				},
				Description: "should do another thing",
				Result:      Passed,
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
	sarifBytes, err := evaluationLog.ToSARIF()
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

	// Check that PhysicalLocation is included in all results
	for _, result := range run.Results {
		require.NotEmpty(t, result.Locations, "Each result should have at least one location")
		for _, location := range result.Locations {
			require.NotNil(t, location.PhysicalLocation, "PhysicalLocation should be present")
			require.Equal(t, informationURI, location.PhysicalLocation.ArtifactLocation.URI, "ArtifactLocation URI should match Metadata.Author.Uri")
			require.Nil(t, location.PhysicalLocation.Region, "Region should be nil for repository-level assessments")
			// LogicalLocations should still be present
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
				Description: "should do a thing",
				Result:      Failed,
				Message:     "thing was not done",
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
	sarifBytes, err := evaluationLog.ToSARIF()
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
