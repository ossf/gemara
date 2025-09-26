package layer4

// EvaluationPlan defines how a set of Layer 4 controls are to be evaluated.
type EvaluationPlan struct {
	Metadata Metadata `json:"metadata" yaml:"metadata"`

	Plans []AssessmentPlan `json:"plans" yaml:"plans"`
}

// Metadata contains metadata about the evaluation plan or evaluation log.
type Metadata struct {
	Id string `json:"id" yaml:"id"`

	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	Evaluator Evaluator `json:"evaluator" yaml:"evaluator"`
}

type Evaluator struct {
	Name string `json:"name" yaml:"name"`

	URI string `json:"uri,omitempty" yaml:"uri,omitempty"`

	Version string `json:"version,omitempty" yaml:"version,omitempty"`

	Contact Contact `json:"contact" yaml:"contact"`
}

type Contact struct {
	// The contact person's name.
	Name string `json:"name" yaml:"name"`

	// Indicates whether this admin is the first point of contact for inquiries. Only one entry should be marked as primary.
	Primary bool `json:"primary" yaml:"primary"`

	// The entity with which the contact is affiliated, such as a school or employer.
	Affiliation *string `json:"affiliation,omitempty" yaml:"affiliation,omitempty"`

	// A preferred email address to reach the contact.
	Email *string `json:"email,omitempty" yaml:"email,omitempty"`

	// A social media handle or profile for the contact.
	Social *string `json:"social,omitempty" yaml:"social,omitempty"`
}

// AssessmentPlan defines all testing procedures for a control id.
type AssessmentPlan struct {
	ControlId string `json:"control-id" yaml:"control-id"`

	Assessments []Assessment `json:"assessments" yaml:"assessments"`
}

// Assessment defines all testing procedures for a requirement.
type Assessment struct {
	// RequirementId is the unique identifier for the requirement being tested.
	RequirementId string `json:"requirement-id" yaml:"requirement-id"`

	// Procedures defines possible testing procedures to evaluate the requirement.
	Procedures []AssessmentProcedure `json:"procedures" yaml:"procedures"`
}

// AssessmentProcedure describes a testing procedure for evaluating a Layer 2 control requirement.
type AssessmentProcedure struct {
	// Id uniquely identifies the assessment procedure being executed
	Id string `json:"id" yaml:"id"`

	// Name provides a summary of the procedure
	Name string `json:"name" yaml:"name"`

	// Description provides a detailed explanation of the procedure
	Description string `json:"description" yaml:"description"`

	// Documentation provides a URL to documentation that describes how the assessment procedure evaluates the control requirement
	Documentation string `json:"documentation,omitempty" yaml:"documentation,omitempty"`
}
