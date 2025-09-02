package layer4

type EvaluationPlan struct {
	Metadata Metadata `json:"metadata"`

	Plans []AssessmentPlan `json:"plans"`
}

type Metadata struct {
	Id string `json:"id"`

	Version string `json:"version,omitempty"`

	Author Contact `json:"author"`
}

type Contact struct {
	// The contact person's name.
	Name string `json:"name"`

	// Indicates whether this admin is the first point of contact for inquiries. Only one entry should be marked as primary.
	Primary bool `json:"primary"`

	// The entity with which the contact is affiliated, such as a school or employer.
	Affiliation *string `json:"affiliation,omitempty"`

	// A preferred email address to reach the contact.
	Email *string `json:"email,omitempty"`

	// A social media handle or profile for the contact.
	Social *string `json:"social,omitempty"`
}

// AssessmentPlan defines all testing procedures for a control id.
type AssessmentPlan struct {
	ControlId string `json:"control-id"`

	Assessments []Assessment `json:"assessments"`
}

// AssessmentPlan defines all testing procedures for a requirement.
type Assessment struct {
	// RequirementID is the unique identifier for the requirement being tested
	RequirementId string `json:"requirement-id"`

	// Procedures defined possible testing procedures to evaluate the requirement.
	Procedures []AssessmentProcedure `json:"procedures"`
}

// AssessmentProcedure describes a testing procedure for evaluating a Layer 2 control requirement.
type AssessmentProcedure struct {
	// Id uniquely identifies the assessment procedure being executed
	Id string `json:"id"`

	// Name provides a summary of the procedure
	Name string `json:"name"`

	// Description provides a detailed explanation of the procedure
	Description string `json:"description"`

	// Documentation provides a URL to documentation that describes how the assessment procedure evaluates the control requirement
	Documentation string `json:"documentation,omitempty"`
}
