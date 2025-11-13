package layer4

import (
	"encoding/json"
	"testing"

	"github.com/ossf/gemara/layer2"
	"github.com/stretchr/testify/require"
)

func TestToSARIF(t *testing.T) {
	testCatalog := makeCatalog("CTRL-1", "Test Control Title", "Test control objective", "REQ-1", "This is the requirement text that should appear in SARIF", "This is the catalog recommendation")

	tests := []struct {
		name          string
		artifactURI   string
		catalog       *layer2.Catalog
		evaluationLog EvaluationLog
		wantRules     int
		wantResults   int
		wantLevels    map[string]string
		wantToolName  string
		wantToolURI   string
		wantToolVer   string
		checkLocation func(*testing.T, *Location)
		checkRule     func(*testing.T, *ReportingDescriptor)
	}{
		{
			name:        "basic conversion with multiple results",
			artifactURI: "",
			catalog:     nil,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "gemara",
				Uri:     "https://github.com/ossf/gemara",
				Version: "1.0.0",
			}, []*AssessmentLog{
				makeAssessmentLog("REQ-1", "should do a thing", Failed, "thing was not done", nil),
				makeAssessmentLog("REQ-2", "should maybe do a thing", NeedsReview, "", nil),
				makeAssessmentLog("REQ-3", "should do another thing", Passed, "", nil),
			}),
			wantRules:   3,
			wantResults: 3,
			wantLevels: map[string]string{
				"REQ-1": "error",
				"REQ-2": "warning",
				"REQ-3": "note",
			},
			wantToolName: "gemara",
			wantToolURI:  "https://github.com/ossf/gemara",
			wantToolVer:  "1.0.0",
			checkLocation: func(t *testing.T, loc *Location) {
				require.NotNil(t, loc.PhysicalLocation)
				require.Equal(t, emptyArtifactURIMessage, loc.PhysicalLocation.ArtifactLocation.URI)
				require.NotEmpty(t, loc.LogicalLocations)
			},
		},
		{
			name:        "with artifactURI parameter",
			artifactURI: "README.md",
			catalog:     nil,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "gemara",
				Uri:     "https://github.com/test/repo",
				Version: "1.0.0",
			}, []*AssessmentLog{
				makeAssessmentLog("REQ-1", "Test requirement", Failed, "Test message", nil),
			}),
			wantRules:   1,
			wantResults: 1,
			wantLevels: map[string]string{
				"REQ-1": "error",
			},
			wantToolName: "gemara",
			wantToolURI:  "https://github.com/test/repo",
			wantToolVer:  "1.0.0",
			checkLocation: func(t *testing.T, loc *Location) {
				require.NotNil(t, loc.PhysicalLocation)
				require.Equal(t, "README.md", loc.PhysicalLocation.ArtifactLocation.URI)
				require.NotEmpty(t, loc.LogicalLocations)
			},
		},
		{
			name:        "empty author URI",
			artifactURI: "",
			catalog:     nil,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "gemara",
				Uri:     "",
				Version: "1.0.0",
			}, []*AssessmentLog{
				makeAssessmentLog("REQ-1", "should do a thing", Failed, "thing was not done", nil),
			}),
			wantRules:   1,
			wantResults: 1,
			wantLevels: map[string]string{
				"REQ-1": "error",
			},
			wantToolName: "gemara",
			wantToolURI:  "",
			wantToolVer:  "1.0.0",
			checkLocation: func(t *testing.T, loc *Location) {
				require.NotNil(t, loc.PhysicalLocation)
				require.Equal(t, emptyArtifactURIMessage, loc.PhysicalLocation.ArtifactLocation.URI)
				require.NotEmpty(t, loc.LogicalLocations)
			},
		},
		{
			name:        "with catalog enrichment",
			artifactURI: "README.md",
			catalog:     testCatalog,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "test-tool",
				Uri:     "https://github.com/test/tool",
				Version: "1.0.0",
			}, []*AssessmentLog{
				{
					Requirement:    Mapping{EntryId: "REQ-1"},
					Description:    "Test description",
					Result:         Failed,
					Message:        "Test failed",
					Recommendation: "Fix this issue by doing X",
					Steps:          []AssessmentStep{func(interface{}) (Result, string) { return Failed, "" }},
					StepsExecuted:  1,
				},
			}),
			wantRules:   1,
			wantResults: 1,
			wantLevels: map[string]string{
				"REQ-1": "error",
			},
			wantToolName: "test-tool",
			wantToolURI:  "https://github.com/test/tool",
			wantToolVer:  "1.0.0",
			checkRule: func(t *testing.T, rule *ReportingDescriptor) {
				require.Equal(t, "REQ-1", rule.ID)
				require.NotNil(t, rule.ShortDescription)
				require.Equal(t, "This is the requirement text that should appear in SARIF", rule.ShortDescription.Text)
				require.NotNil(t, rule.FullDescription)
				require.Contains(t, rule.FullDescription.Text, "Test control objective")
				require.Contains(t, rule.FullDescription.Text, "This is the requirement text")
				require.NotNil(t, rule.Help)
				require.Equal(t, "Fix this issue by doing X", rule.Help.Text, "should prefer AssessmentLog recommendation over catalog")
				require.Empty(t, rule.HelpUri)
			},
		},
		{
			name:        "without catalog",
			artifactURI: "README.md",
			catalog:     nil,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "test-tool",
				Uri:     "https://github.com/test/tool",
				Version: "1.0.0",
			}, []*AssessmentLog{
				{
					Requirement:    Mapping{EntryId: "REQ-1"},
					Description:    "Test description",
					Result:         Failed,
					Message:        "Test failed",
					Recommendation: "Fix this issue by doing X",
					Steps:          []AssessmentStep{func(interface{}) (Result, string) { return Failed, "" }},
					StepsExecuted:  1,
				},
			}),
			wantRules:   1,
			wantResults: 1,
			wantLevels: map[string]string{
				"REQ-1": "error",
			},
			wantToolName: "test-tool",
			wantToolURI:  "https://github.com/test/tool",
			wantToolVer:  "1.0.0",
			checkRule: func(t *testing.T, rule *ReportingDescriptor) {
				require.Equal(t, "REQ-1", rule.ID)
				require.Nil(t, rule.ShortDescription)
				require.Nil(t, rule.FullDescription)
				require.Nil(t, rule.Help)
				require.Empty(t, rule.HelpUri)
			},
		},
		{
			name:        "catalog recommendation when assessment log has none",
			artifactURI: "README.md",
			catalog:     testCatalog,
			evaluationLog: makeEvaluationLog(Author{
				Name:    "test-tool",
				Uri:     "https://github.com/test/tool",
				Version: "1.0.0",
			}, []*AssessmentLog{
				{
					Requirement:   Mapping{EntryId: "REQ-1"},
					Description:   "Test description",
					Result:        Failed,
					Message:       "Test failed",
					Steps:         []AssessmentStep{func(interface{}) (Result, string) { return Failed, "" }},
					StepsExecuted: 1,
				},
			}),
			wantRules:   1,
			wantResults: 1,
			wantLevels: map[string]string{
				"REQ-1": "error",
			},
			wantToolName: "test-tool",
			wantToolURI:  "https://github.com/test/tool",
			wantToolVer:  "1.0.0",
			checkRule: func(t *testing.T, rule *ReportingDescriptor) {
				require.NotNil(t, rule.Help)
				require.Equal(t, "This is the catalog recommendation", rule.Help.Text, "should use catalog recommendation when assessment log has none")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sarifBytes, err := tt.evaluationLog.ToSARIF(tt.artifactURI, tt.catalog)
			require.NoError(t, err)

			sarif := toSARIFReport(t, sarifBytes)
			require.Len(t, sarif.Runs, 1)

			run := sarif.Runs[0]

			require.Len(t, run.Tool.Driver.Rules, tt.wantRules)
			require.Len(t, run.Results, tt.wantResults)

			require.Equal(t, tt.wantToolName, run.Tool.Driver.Name)
			require.Equal(t, tt.wantToolURI, run.Tool.Driver.InformationURI)
			require.Equal(t, tt.wantToolVer, run.Tool.Driver.Version)

			levels := make(map[string]string)
			for _, r := range run.Results {
				levels[r.RuleID] = r.Level
				if tt.checkLocation != nil {
					require.NotEmpty(t, r.Locations)
					tt.checkLocation(t, &r.Locations[0])
				}
			}

			for ruleID, wantLevel := range tt.wantLevels {
				require.Equal(t, wantLevel, levels[ruleID], "rule %s should have level %s", ruleID, wantLevel)
			}

			if tt.checkRule != nil && len(run.Tool.Driver.Rules) > 0 {
				tt.checkRule(t, &run.Tool.Driver.Rules[0])
			}

			_, err = json.Marshal(sarif)
			require.NoError(t, err)
		})
	}
}

func TestToSARIF_ResultLevels(t *testing.T) {
	tests := []struct {
		result    Result
		wantLevel string
		wantCount int
	}{
		{Failed, "error", 1},
		{NeedsReview, "warning", 1},
		{Unknown, "warning", 1},
		{Passed, "note", 1},
		{NotApplicable, "", 0},
		{NotRun, "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.result.String(), func(t *testing.T) {
			evaluationLog := makeEvaluationLog(Author{
				Name:    "test",
				Uri:     "https://test",
				Version: "1.0.0",
			}, []*AssessmentLog{
				makeAssessmentLog("REQ-1", "test", tt.result, "", nil),
			})

			sarifBytes, err := evaluationLog.ToSARIF("", nil)
			require.NoError(t, err)

			sarif := toSARIFReport(t, sarifBytes)
			require.Len(t, sarif.Runs[0].Results, tt.wantCount)

			if tt.wantCount > 0 {
				require.Equal(t, tt.wantLevel, sarif.Runs[0].Results[0].Level)
			}
		})
	}
}

// Helper functions

func makeEvaluationLog(author Author, logs []*AssessmentLog) EvaluationLog {
	return EvaluationLog{
		Evaluations: []*ControlEvaluation{
			{
				Name:           "Example Control",
				Control:        Mapping{EntryId: "CTRL-1"},
				Result:         Passed,
				AssessmentLogs: logs,
			},
		},
		Metadata: Metadata{Author: author},
	}
}

func makeAssessmentLog(entryID, description string, result Result, message string, steps []AssessmentStep) *AssessmentLog {
	if steps == nil {
		steps = []AssessmentStep{func(interface{}) (Result, string) { return result, "" }}
	}
	return &AssessmentLog{
		Requirement:   Mapping{EntryId: entryID},
		Description:   description,
		Result:        result,
		Message:       message,
		Steps:         steps,
		StepsExecuted: int64(len(steps)),
	}
}

func makeCatalog(controlID, controlTitle, controlObjective, reqID, reqText, reqRecommendation string) *layer2.Catalog {
	return &layer2.Catalog{
		ControlFamilies: []layer2.ControlFamily{
			{
				Id:    "test-family",
				Title: "Test Family",
				Controls: []layer2.Control{
					{
						Id:        controlID,
						Title:     controlTitle,
						Objective: controlObjective,
						AssessmentRequirements: []layer2.AssessmentRequirement{
							{
								Id:             reqID,
								Text:           reqText,
								Recommendation: reqRecommendation,
							},
						},
					},
				},
			},
		},
	}
}

func toSARIFReport(t *testing.T, data []byte) *SarifReport {
	t.Helper()
	var sarif SarifReport
	err := json.Unmarshal(data, &sarif)
	require.NoError(t, err)
	require.NotNil(t, &sarif)
	return &sarif
}
