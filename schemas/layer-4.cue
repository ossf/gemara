package schemas

import "time"

@go(layer4)

// EvaluationDocument defines how a set of Layer 2 controls are to be evaluated and the associated outcomes of the evaluation.
#EvaluationDocument: {
	metadata: #Metadata
	plan?:    #EvaluationPlan @go(Plan,optional=nillable)
	logs?: [...#EvaluationLog]
}

// Metadata contains metadata about the Layer 4 evaluation document.
#Metadata: {
	id:       string
	version?: string
	author:   #Author
	"mapping-references"?: [...#MappingReference] @go(MappingReferences) @yaml("mapping-references,omitempty")
}

// Author contains the information about the entity that produced the evaluation document.
#Author: {
	name:     string
	uri?:     string
	version?: string
	contact?: #Contact @go(Contact)
}

// EvaluationPlan defines how a set of Layer 2 controls are to be evaluated.
#EvaluationPlan: {
	plans: [...#AssessmentPlan]
	// Executors defines the assessment executors that can be used to execute assessment procedures.
	executors?: [...#AssessmentExecutor]
}

// EvaluationLog contains the results of evaluating a set of Layer 2 controls.
#EvaluationLog: {
	// Executor defines the assessment executor that produces the log
	executor?: #AssessmentExecutor @go(Executor,option=nillable)
	evaluations: [#ControlEvaluation, ...#ControlEvaluation] @go(Evaluations,type=[]*ControlEvaluation)
}

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
	// Strategy defines the rules for aggregating results from multiple procedures when they conflict.
	strategy?: #Strategy @go(Strategy,optional=nillable)
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
	// Executors lists which assessment executors can execute this procedure.
	executors: [...#ExecutorMapping]
	// Strategy defines the rules for aggregating results from multiple executors running the same procedure.
	strategy?: #Strategy @go(Strategy,optional=nillable)
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
	// Confidence describes the confidence level in this assessment result on a scale of 1 to 10, with 10 being the most confident.
	confidence?: int & >=1 & <=10
}

// ExecutorMapping maps an assessment executor to a procedure.
#ExecutorMapping: {
	id: string
	// Role determines how this executor participates in conflict resolution when using AdvisoryRequiresConfirmation strategy.
	role?: #ExecutorRole @go(Role)
	// Remarks provides context about why this executor-procedure combination was chosen.
	remarks?: string
}

// AssessmentExecutor describes an assessment method (tool or manual approach) that can be used to execute assessment procedures.
#AssessmentExecutor: {
	// Id uniquely identifies the assessment method.
	id: string
	// Name provides the name of the assessment method.
	name: string
	// Type specifies whether the executor is automated or manual.
	type: #ExecutorType @go(Type)
	// Version specifies the version of the assessment method (if applicable, e.g., for tools).
	version?: string
	// Description provides additional context about the assessment method.
	description?: string
	// Documentation provides a URL to documentation for the assessment method.
	documentation?: =~"^https?://[^\\s]+$"
}

// Strategy defines the rules for resolving conflicts between multiple executor results.
#Strategy: {
	// ConflictRuleType specifies the aggregation logic used to resolve conflicts when multiple executors provide results for the same assessment procedure.
	"conflict-rule-type": #ConflictRuleType @go(ConflictRuleType)
	// Remarks provides context for why this specific conflict resolution strategy was chosen.
	remarks?: string
}

// MappingReference references an external mapping document.
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

#AssessmentStep: string @go(-)

#Result: "Not Run" | "Passed" | "Failed" | "Needs Review" | "Not Applicable" | "Unknown" @go(-)

#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime)

// ConflictRuleType defines how to resolve conflicts when multiple executors produce different results.
// - Strict: if any executor reports a failure, the overall procedure result is failed.
// - ManualOverride: manual review executors take precedence over automated executors when results conflict.
// - AdvisoryRequiresConfirmation: Advisory executors require confirmation from Primary executors before triggering findings.
#ConflictRuleType: "Strict" | "ManualOverride" | "AdvisoryRequiresConfirmation" @go(-)

// ExecutorType specifies whether the executor is automated or manual.
// Automated executors are tools or scripts that run without human intervention.
// Manual executors require human review or judgment.
#ExecutorType: "Automated" | "Manual" @go(-)

// ExecutorRole determines how this executor participates in conflict resolution when using AdvisoryRequiresConfirmation strategy.
#ExecutorRole: "Primary" | "Advisory" @go(-)

#Email: =~"^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"
