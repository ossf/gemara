package gemara

// This file is for reusable test data to help seed ideas and reduce duplication.

var (
	// Generic applicability for testing
	testingApplicability = []string{"test-applicability"}

	// Assessment Results
	passingAssessmentStep = func(interface{}) (Result, string) {
		return Passed, ""
	}
	failingAssessmentStep = func(interface{}) (Result, string) {
		return Failed, ""
	}
	needsReviewAssessmentStep = func(interface{}) (Result, string) {
		return NeedsReview, ""
	}
	unknownAssessmentStep = func(interface{}) (Result, string) {
		return Unknown, ""
	}
)

func failingAssessmentPtr() *AssessmentLog {
	a := failingAssessment()
	return &a
}

func failingAssessment() AssessmentLog {
	return AssessmentLog{
		Requirement: EntryMapping{
			EntryId: "failingAssessment()",
		},
		Description: "failing assessment",
		Steps: []AssessmentStep{
			failingAssessmentStep,
			passingAssessmentStep,
		},
		Applicability: testingApplicability,
	}
}
func passingAssessmentPtr() *AssessmentLog {
	a := passingAssessment()
	return &a
}

func passingAssessment() AssessmentLog {
	return AssessmentLog{
		Requirement: EntryMapping{
			EntryId: "passingAssessment()",
		},
		Description: "passing assessment",
		Steps: []AssessmentStep{
			passingAssessmentStep,
		},
		Applicability: testingApplicability,
	}
}
func needsReviewAssessmentPtr() *AssessmentLog {
	a := needsReviewAssessment()
	return &a
}

func needsReviewAssessment() AssessmentLog {
	return AssessmentLog{
		Requirement: EntryMapping{
			EntryId: "needsReviewAssessment()",
		},
		Description: "needs review assessment",
		Steps: []AssessmentStep{
			passingAssessmentStep,
			needsReviewAssessmentStep,
			passingAssessmentStep,
		},
		Applicability: testingApplicability,
	}
}
func unknownAssessmentPtr() *AssessmentLog {
	a := unknownAssessment()
	return &a
}

func unknownAssessment() AssessmentLog {
	return AssessmentLog{
		Requirement: EntryMapping{
			EntryId: "unknownAssessment()",
		},
		Description: "unknown assessment",
		Steps: []AssessmentStep{
			passingAssessmentStep,
			unknownAssessmentStep,
			passingAssessmentStep,
		},
		Applicability: testingApplicability,
	}
}

func badRevertPassingAssessment() AssessmentLog {
	return AssessmentLog{
		Requirement: EntryMapping{
			EntryId: "badRevertPassingAssessment()",
		},
		Description: "bad revert passing assessment",
		Steps: []AssessmentStep{
			passingAssessmentStep,
			passingAssessmentStep,
			passingAssessmentStep,
			passingAssessmentStep,
		},
		Applicability: testingApplicability,
	}
}
