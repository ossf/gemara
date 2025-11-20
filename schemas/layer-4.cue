package schemas

import "time"

@go(layer4)

// EvaluationDocument defines how a set of Layer 2 controls are to be evaluated and the associated outcomes of the evaluation.
#EvaluationDocument: {
	metadata:           #Metadata
	"evaluation-plan"?: #EvaluationPlan @go(EvaluationPlan,optional=nillable)
	"evaluation-logs"?: [...#EvaluationLog] @go(EvaluationLogs)
}

// EvaluationPlan defines how a set of Layer 2 controls are to be evaluated.
#EvaluationPlan: {
	metadata: #Metadata
	// Evaluators defines the assessment evaluators that can be used to execute assessment procedures.
	evaluators: [...#Actor] @go(Evaluators)
	plans: [...#AssessmentPlan]
}

// EvaluationLog contains the results of evaluating a set of Layer 2 controls.
#EvaluationLog: {
	"metadata"?: #Metadata
	"evaluations": [#ControlEvaluation, ...#ControlEvaluation] @go(Evaluations,type=[]*ControlEvaluation)
}

// Metadata contains common fields shared across all metadata types.
#Metadata: {
	id:           string
	version?:     string
	author:       #Actor
	description?: string
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

// Actor represents an entity (human or tool) that can perform actions in evaluations.
#Actor: {
	// Id uniquely identifies the actor.
	id: string
	// Name provides the name of the actor.
	name: string
	// Type specifies how the evaluation is executed (automated tool or manual/human evaluator).
	type: #ExecutionType @go(Type)
	// Version specifies the version of the actor (if applicable, e.g., for tools).
	version?: string
	// Description provides additional context about the actor.
	description?: string
	// Uri provides a general URI for the actor information.
	uri?: =~"^https?://[^\\s]+$"
	// Contact provides contact information for the actor.
	contact?: #Contact @go(Contact)
}

// ExecutionType specifies how the evaluation is executed (automated or manual).
#ExecutionType: "Automated" | "Manual" @go(-)

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
	// ConfidenceLevel indicates the evaluator's confidence level in this specific assessment result.
	"confidence-level"?: #ConfidenceLevel @go(ConfidenceLevel)
}

#AssessmentStep: string @go(-)

#Result: "Not Run" | "Passed" | "Failed" | "Needs Review" | "Not Applicable" | "Unknown" @go(-)

#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime)

// ConfidenceLevel indicates the evaluator's confidence level in an assessment result.
#ConfidenceLevel: "Undetermined" | "Low" | "Medium" | "High" @go(-)

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
	// Evaluators lists which assessment evaluators can execute this procedure.
	evaluators: [...#EvaluatorMapping] @go(Evaluators)
	// Strategy defines the rules for aggregating results from multiple evaluators running the same procedure.
	strategy?: #Strategy @go(Strategy,optional=nillable)
}

// EvaluatorMapping maps an assessment evaluator to a procedure.
#EvaluatorMapping: {
	id: string
	// Authoritative determines how this evaluator participates in conflict resolution when using AuthoritativeConfirmation strategy.
	// If true, the evaluator can trigger findings independently. If false (default), the evaluator requires confirmation from authoritative evaluators before triggering findings.
	// Note: This field is only used with AuthoritativeConfirmation strategy. With Strict (default) or ManualOverride strategies, this field is ignored.
	authoritative?: bool @go(Authoritative)
	// Remarks provides context about why this evaluator-procedure combination was chosen.
	remarks?: string
}

// Strategy defines the rules for resolving conflicts between multiple evaluator results.
#Strategy: {
	// ConflictRuleType specifies the aggregation logic used to resolve conflicts when multiple evaluators provide results for the same assessment procedure.
	"conflict-rule-type": #ConflictRuleType @go(ConflictRuleType)
	// Remarks provides context for why this specific conflict resolution strategy was chosen.
	remarks?: string
}

// ConflictRuleType defines how to resolve conflicts when multiple evaluators produce different results.
#ConflictRuleType: "Strict" | "ManualOverride" | "AuthoritativeConfirmation" @go(-)

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
