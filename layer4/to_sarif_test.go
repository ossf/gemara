package layer4

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToSARIF(t *testing.T) {
	var sarif *SarifReport
	ce := &ControlEvaluation{
		Name:      "Example Control",
		ControlID: "CTRL-1",
		Result:    Passed,
		AssessmentLogs: []*AssessmentLog{
			{
				RequirementId: "REQ-1",
				Description:   "should do a thing",
				Result:        Failed,
				Message:       "thing was not done",
			},
			{
				RequirementId: "REQ-2",
				Description:   "should maybe do a thing",
				Result:        NeedsReview,
			},
			{
				RequirementId: "REQ-3",
				Description:   "should do another thing",
				Result:        Passed,
			},
		},
	}
	informationURI := "https://github.com/ossf/gemara"
	version := "1.0.0"

	evaluationLog := EvaluationLog{
		Evaluations: []*ControlEvaluation{ce},
		Metadata: Metadata{
			Evaluator: Evaluator{
				Name:    "gemara",
				URI:     informationURI,
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

	// ensure JSON marshals cleanly
	_, err = json.Marshal(sarif)
	require.NoError(t, err)
}
