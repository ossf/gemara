package schemas

import "time"

// EvaluationPlan defines how a set of Layer 2 controls are to be evaluated.
#EvaluationPlan: {
	metadata: #Metadata
	plans: [...#AssessmentPlan]
}

// EvaluationLog contains the results of evaluating a set of Layer 2 controls.
#EvaluationLog: {
	"evaluations": [#ControlEvaluation, ...#ControlEvaluation] @go(Evaluations)
	"metadata"?: #Metadata @go(Metadata)
}

// Metadata contains metadata about the Layer 4 evaluation plan and log.
#Metadata: {
	id:        string
	version?:  string
	evaluator: #Evaluator
}

// Evaluator contains the information about the entity that produced the evaluation results.
#Evaluator: {
	"name":     string
	"uri"?:     string
	"version"?: string
	"contact"?: #Contact @go(Contact)
}

// ControlEvaluation contains the results of evaluating a single Layer 4 control.
// TODO: are all control requirements guaranteed to be evaluated at once, or does applicability influence this?
#ControlEvaluation: {
	name:              string
	"control-id":      string @go(ControlId)
	result:            #Result
	message:           string
	"corrupted-state": bool @go(CorruptedState)
	"assessment-logs": [...#AssessmentLog] @go(AssessmentLogs)
}

// AssessmentLog contains the results of executing a single assessment procedure for a control requirement.
#AssessmentLog: {
	// RequirementId identifies the control requirement assessed.
	"requirement-id": string @go(RequirementId)
	// ProcedureId uniquely identifies the assessment procedure associated with the log
	"procedure-id"?: string @go(ProcedureId)
	applicability: [...string]
	description: string
	result:      #Result
	message:     string
	steps: [...#AssessmentStep]
	"steps-executed"?: int @go(StepsExecuted)
	"start":           #Datetime
	"end"?:            #Datetime
	value?:            _
	changes?: {[string]: #Change}
	recommendation?: string
}

#AssessmentStep: string

#Change: {
	"target-name":    string @go(TargetName)
	description:      string
	"target-object"?: _ @go(TargetObject)
	applied?:         bool
	reverted?:        bool
	error?:           string
}

#Result: "Not Run" | "Passed" | "Failed" | "Needs Review" | "Not Applicable" | "Unknown"

#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime,format="date-time")

// AssessmentPlan defines all testing procedures for a control id.
#AssessmentPlan: {
	"control-id": string @go(ControlId)
	"assessments": [...#Assessment] @go(Assessments)
}

// Assessment defines all testing procedures for a requirement.
#Assessment: {
	// RequirementId is the unique identifier for the requirement being tested.
	"requirement-id": string @go(RequirementId)
	// Procedures defines possible testing procedures to evaluate the requirement.
	procedures: [...#AssessmentProcedure] @go(Procedures)
}

// AssessmentProcedure describes a testing procedure for evaluating a Layer 2 control requirement.
#AssessmentProcedure: {
	// Id uniquely identifies the assessment procedure being executed
	id: string
	// Name provides a summary of the procedure
	name: string
	// Description provides a detailed explanation of the procedure
	description: string
	// Documentation provides a URL to documentation that describes how the assessment procedure evaluates the control requirement
	documentation?: =~"^https?://[^\\s]+$"
}

#Contact: {
	// The contact person's name.
	name: string
	// Indicates whether this admin is the first point of contact for inquiries. Only one entry should be marked as primary.
	primary: bool
	// The entity with which the contact is affiliated, such as a school or employer.
	affiliation?: string @go(Affiliation,type=*string)
	// A preferred email address to reach the contact.
	email?: #Email @go(Email,type=*Email)
	// A social media handle or profile for the contact.
	social?: string @go(Social,type=*string)
}

#Email: =~"^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"
