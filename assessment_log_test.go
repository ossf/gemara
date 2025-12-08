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
			aggregator := NewConfidenceAggregator()
			result := anyOldAssessment.runStep(nil, test.step, aggregator)
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
		expectedConfidence ConfidenceLevel
	}{
		{
			name: "Result changed - step determines final result",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Failed, "Step 2 failed", Low
				},
			},
			expectedConfidence: Medium, // 50% High, 50% Low → 50% Medium+ threshold met
		},
		{
			name: "Result changed - Passed to NeedsReview",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 2 needs review", Medium
				},
			},
			expectedConfidence: Medium,
		},
		{
			name: "Same result - multiple Passed steps",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", Medium
				},
			},
			expectedConfidence: Medium, // Minimum of High and Medium
		},
		{
			name: "Same result - multiple Passed steps with Low confidence",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", Medium
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", Low
				},
			},
			expectedConfidence: Medium, // 50% Medium, 50% Low → 50% Medium+ threshold met
		},
		{
			name: "Result did not change - Failed persists despite Passed step",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Failed, "Step 1 failed", Low
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", High
				},
			},
			expectedConfidence: Low, // Confidence from step1 (step2 didn't affect result)
		},
		{
			name: "Result did not change - Unknown persists despite Passed step",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Unknown, "Step 1 unknown", Undetermined
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", High
				},
			},
			expectedConfidence: Undetermined, // Confidence from step1 (step2 didn't affect result)
		},
		{
			name: "Result did not change - NeedsReview persists despite Passed step",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NeedsReview, "Step 1 needs review", Medium
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 2 passed", High
				},
			},
			expectedConfidence: Medium, // Confidence from step1 (step2 didn't affect result)
		},
		{
			name: "Result did not change - Passed persists despite NotRun step",
			steps: []AssessmentStep{
				func(interface{}) (Result, string, ConfidenceLevel) {
					return Passed, "Step 1 passed", High
				},
				func(interface{}) (Result, string, ConfidenceLevel) {
					return NotRun, "Step 2 not run", Undetermined
				},
			},
			expectedConfidence: High, // Confidence from step1 (step2 didn't affect result)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assessment, err := NewAssessment("test-id", "test description", []string{"test"}, tt.steps)
			require.NoError(t, err)
			assessment.Run(nil)

			assert.Equal(t, tt.expectedConfidence, assessment.ConfidenceLevel)
		})
	}
}
