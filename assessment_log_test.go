package gemara

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getAssessmentsTestData() []struct {
	testName           string
	assessment         AssessmentLog
	numberOfSteps      int
	numberOfStepsToRun int
	expectedResult     Result
} {
	return []struct {
		testName           string
		assessment         AssessmentLog
		numberOfSteps      int
		numberOfStepsToRun int
		expectedResult     Result
	}{
		{
			testName:   "AssessmentLog with no steps",
			assessment: AssessmentLog{},
		},
		{
			testName:           "AssessmentLog with one step",
			assessment:         passingAssessment(),
			numberOfSteps:      1,
			numberOfStepsToRun: 1,
			expectedResult:     Passed,
		},
		{
			testName:           "AssessmentLog with two steps",
			assessment:         failingAssessment(),
			numberOfSteps:      2,
			numberOfStepsToRun: 1,
			expectedResult:     Failed,
		},
		{
			testName:           "AssessmentLog with three steps",
			assessment:         needsReviewAssessment(),
			numberOfSteps:      3,
			numberOfStepsToRun: 3,
			expectedResult:     NeedsReview,
		},
		{
			testName:           "AssessmentLog with four steps",
			assessment:         badRevertPassingAssessment(),
			numberOfSteps:      4,
			numberOfStepsToRun: 4,
			expectedResult:     Passed,
		},
	}
}

// TestNewStep ensures that NewStep queues a new step in the AssessmentLog
func TestAddStep(t *testing.T) {
	for _, test := range getAssessmentsTestData() {
		t.Run(test.testName, func(t *testing.T) {
			if len(test.assessment.Steps) != test.numberOfSteps {
				t.Errorf("Bad test data: expected to start with %d, got %d", test.numberOfSteps, len(test.assessment.Steps))
			}
			test.assessment.AddStep(passingAssessmentStep)
			if len(test.assessment.Steps) != test.numberOfSteps+1 {
				t.Errorf("expected %d, got %d", test.numberOfSteps, len(test.assessment.Steps))
			}
		})
	}
}

// TestRunStep ensures that runStep runs the step and updates the AssessmentLog
func TestRunStep(t *testing.T) {
	stepsTestData := []struct {
		testName        string
		step            AssessmentStep
		result          Result
		confidenceLevel ConfidenceLevel
	}{
		{
			testName:        "Failing step",
			step:            failingAssessmentStep,
			result:          Failed,
			confidenceLevel: Low,
		},
		{
			testName:        "Passing step",
			step:            passingAssessmentStep,
			result:          Passed,
			confidenceLevel: High,
		},
		{
			testName:        "Needs review step",
			step:            needsReviewAssessmentStep,
			result:          NeedsReview,
			confidenceLevel: Medium,
		},
		{
			testName:        "Unknown step",
			step:            unknownAssessmentStep,
			result:          Unknown,
			confidenceLevel: Undetermined,
		},
	}
	for _, test := range stepsTestData {
		t.Run(test.testName, func(t *testing.T) {
			anyOldAssessment := AssessmentLog{}
			result := anyOldAssessment.runStep(nil, test.step)
			if result != test.result {
				t.Errorf("expected %s, got %s", test.result, result)
			}
			if anyOldAssessment.Result != test.result {
				t.Errorf("expected %s, got %s", test.result, anyOldAssessment.Result)
			}
			if anyOldAssessment.ConfidenceLevel != test.confidenceLevel {
				t.Errorf("expected confidence %s, got %s", test.confidenceLevel, anyOldAssessment.ConfidenceLevel)
			}
		})
	}
}

// TestRun ensures that Run executes all steps, halting if any step does not return Passed
func TestRun(t *testing.T) {
	for _, data := range getAssessmentsTestData() {
		t.Run(data.testName, func(t *testing.T) {
			a := data.assessment // copy the assessment to prevent duplicate executions in the next test
			result := a.Run(nil)
			if result != a.Result {
				t.Errorf("expected match between Run return value (%s) and assessment Result value (%s)", result, data.expectedResult)
			}
			if a.StepsExecuted != int64(data.numberOfStepsToRun) {
				t.Errorf("expected to run %d tests, got %d", data.numberOfStepsToRun, a.StepsExecuted)
			}
		})
	}
}

func TestNewAssessment(t *testing.T) {
	newAssessmentsTestData := []struct {
		testName      string
		requirementId string
		description   string
		applicability []string
		steps         []AssessmentStep
		expectedError bool
	}{
		{
			testName:      "Empty requirementId",
			requirementId: "",
			description:   "test",
			applicability: []string{"test"},
			steps:         []AssessmentStep{passingAssessmentStep},
			expectedError: true,
		},
		{
			testName:      "Empty description",
			requirementId: "test",
			description:   "",
			applicability: []string{"test"},
			steps:         []AssessmentStep{passingAssessmentStep},
			expectedError: true,
		},
		{
			testName:      "Empty applicability",
			requirementId: "test",
			description:   "test",
			applicability: []string{},
			steps:         []AssessmentStep{passingAssessmentStep},
			expectedError: true,
		},
		{
			testName:      "Empty steps",
			requirementId: "test",
			description:   "test",
			applicability: []string{"test"},
			steps:         []AssessmentStep{},
			expectedError: true,
		},
		{
			testName:      "Good data",
			requirementId: "test",
			description:   "test",
			applicability: []string{"test"},
			steps:         []AssessmentStep{passingAssessmentStep},
			expectedError: false,
		},
	}
	for _, data := range newAssessmentsTestData {
		t.Run(data.testName, func(t *testing.T) {
			assessment, err := NewAssessment(data.requirementId, data.description, data.applicability, data.steps)
			if data.expectedError && err == nil {
				t.Error("expected error, got nil")
			}
			if !data.expectedError && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if assessment == nil && !data.expectedError {
				t.Error("expected assessment object, got nil")
			}
		})
	}
}

func TestConfidenceLevelFromSteps(t *testing.T) {
	tests := []struct {
		name               string
		steps              []AssessmentStep
		expectedResult     Result
		expectedConfidence ConfidenceLevel
	}{
		{
			name: "Passed then NeedsReview",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 2 needs review", Medium
				},
			},
			expectedResult:     NeedsReview,
			expectedConfidence: Medium,
		},
		{
			name: "Passed then Passed with different confidence",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", Low
				},
			},
			expectedResult:     Passed,
			expectedConfidence: Low, // Use last step's confidence
		},
		{
			name: "NeedsReview then Passed",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 1 needs review", Medium
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", High
				},
			},
			expectedResult:     NeedsReview,
			expectedConfidence: High, // Use last step's confidence
		},
		{
			name: "Multiple NeedsReview steps",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 1 needs review", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 2 needs review", Low
				},
			},
			expectedResult:     NeedsReview,
			expectedConfidence: Low,
		},
		{
			name: "Unknown then Passed",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Unknown, "Step 1 unknown", Undetermined
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", High
				},
			},
			expectedResult:     Unknown,
			expectedConfidence: High, // Use last step's confidence
		},
		{
			name: "Passed then Unknown",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Unknown, "Step 2 unknown", Undetermined
				},
			},
			expectedResult:     Unknown,
			expectedConfidence: Undetermined,
		},
		{
			name: "Failed stops execution",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Failed, "Step 2 failed", Low
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 3 passed", High
				},
			},
			expectedResult:     Failed,
			expectedConfidence: Low,
		},
		{
			name: "Prereqs, then NeedsReview, then Passed",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Prereq 1", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Prereq 2", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Check 1", Medium
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Check 2", High
				},
			},
			expectedResult:     NeedsReview,
			expectedConfidence: High, // Use last step's confidence
		},
		{
			name: "Prereqs, then Passed Low, then Passed High",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Prereq 1", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Prereq 2", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Check 1", Low
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Check 2", High
				},
			},
			expectedResult:     Passed,
			expectedConfidence: High, // Use last step's confidence
		},
		{
			name: "Passed High, then NeedsReview Low, then Passed High",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 2", Low
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 3", High
				},
			},
			expectedResult:     NeedsReview,
			expectedConfidence: High, // Use last step's confidence
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment, err := NewAssessment("test-id", "test description", []string{"test"}, tt.steps)
			require.NoError(t, err)
			result := assessment.Run(nil)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedConfidence, assessment.ConfidenceLevel)
		})
	}
}
