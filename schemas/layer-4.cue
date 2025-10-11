package schemas

import "time"

@go(layer4)

// EvaluationPlan defines how a set of Layer 2 controls are to be evaluated.
#EvaluationPlan: {
	metadata: #Metadata
	plans: [...#AssessmentPlan]
}

// EvaluationLog contains the results of evaluating a set of Layer 2 controls.
#EvaluationLog: {
	"evaluations": [#ControlEvaluation, ...#ControlEvaluation] @go(Evaluations,type=[]*ControlEvaluation)
	"metadata"?: #Metadata @go(Metadata)
}

// Metadata contains metadata about the Layer 4 evaluation plan and log.
#Metadata: {
	id:       string
	version?: string
	author:   #Author
	"mapping-references"?: [...#MappingReference] @go(MappingReferences) @yaml("mapping-references,omitempty")
}

#MappingReference: {
	id:           string
	title:        string
	version:      string
	description?: string
	url?:         =~"^https?://[^\\s]+$"
}

#Mapping: {
	// ReferenceId should reference the corresponding MappingReference id
	"reference-id": string @go(ReferenceId)
	// EntryId should reference the specific element within the referenced document
	"entry-id": string @go(EntryId)
	// Strength describes how effectively the referenced item addresses the associated control or procedure on a scale of 1 to 10, with 10 being the most effective.
	strength?: int & >=1 & <=10
	// Remarks provides additional context about the mapping entry.
	remarks?: string
}

// Author contains the information about the entity that produced the evaluation plan or log.
#Author: {
	"name":     string
	"uri"?:     string
	"version"?: string
	"contact"?: #Contact @go(Contact)
}

// ControlEvaluation contains the results of evaluating a single Layer 4 control.
#ControlEvaluation: {
	name:    string
	result:  #Result
	message: string
	control: #Mapping
	"assessment-logs": [...#AssessmentLog] @go(AssessmentLogs,type=[]*AssessmentLog)
	// Enforce that control reference and the assessments' references match
	// This formulation uses the control's reference if the assessment doesn't include a reference
	"assessment-logs": [...{
		requirement: "reference-id": (control."reference-id")
	}] @go(AssessmentLogs,type=[]*AssessmentLog)
}

// AssessmentLog contains the results of executing a single assessment procedure for a control requirement.
#AssessmentLog: {
	// Requirement should map to the assessment requirement for this assessment.
	requirement: #Mapping
	// Procedure should map to the assessment procedure being executed.
	procedure: #Mapping
	// Description provides a summary of the assessment procedure.
	description: string
	// Result is the overall outcome of the assessment procedure, matching the result of the last step that was run.
	result: #Result
	// Message provides additional context about the assessment result.
	message: string
	// Applicability is elevated from the Layer 2 Assessment Requirement to aid in execution and reporting.
	applicability: [...string] @go(Applicability,type=[]string)
	// Steps are sequential actions taken as part of the assessment, which may halt the assessment if a failure occurs.
	steps: [...#AssessmentStep]
	// Steps-executed is the number of steps that were executed as part of the assessment.
	"steps-executed"?: int @go(StepsExecuted)
	// Start is the timestamp when the assessment began.
	start: #Datetime
	// End is the timestamp when the assessment concluded.
	end?: #Datetime
	// Recommendation provides guidance on how to address a failed assessment.
	recommendation?: string
}

#AssessmentStep: string @go(-)

#Result: "Not Run" | "Passed" | "Failed" | "Needs Review" | "Not Applicable" | "Unknown" @go(-)

#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime)

// AssessmentPlan defines all testing procedures for a control id.
#AssessmentPlan: {
	// Control points to the Layer 2 control being evaluated.
	control: #Mapping
	// Assessments defines possible testing procedures to evaluate the control.
	assessments: [...#Assessment] @go(Assessments,type=[]Assessment)
	// Enforce that control reference and the assessments' references match
	// This formulation uses the control's reference if the assessment doesn't include a reference
	assessments: [...{
		requirement: "reference-id": (control."reference-id")
	}] @go(Assessments,type=[]Assessment)
}

// Assessment defines all testing procedures for a requirement.
#Assessment: {
	// RequirementId points to the requirement being tested.
	requirement: #Mapping
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
