package schemas

@go(gemara)

// EvaluationPlan defines how a set of Layer 2 controls are to be evaluated.
#EvaluationPlan: {
	"metadata"?: #Metadata @go(Metadata)
	plans: [...#AssessmentPlan]
}

// EvaluationLog contains the results of evaluating a set of Layer 2 controls.
#EvaluationLog: {
	"evaluations": [#ControlEvaluation, ...#ControlEvaluation] @go(Evaluations,type=[]*ControlEvaluation)
	"metadata"?: #Metadata @go(Metadata)
}

// ControlEvaluation contains the results of evaluating a single Layer 4 control.
#ControlEvaluation: {
	name:    string
	result:  #Result
	message: string
	control: #SingleMapping
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
	requirement: #SingleMapping
	// Procedure should map to the assessment procedure being executed.
	procedure: #SingleMapping
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

// ConfidenceLevel indicates the evaluator's confidence level in an assessment result.
#ConfidenceLevel: "Undetermined" | "Low" | "Medium" | "High" @go(-)

// AssessmentPlan defines all testing procedures for a control id.
#AssessmentPlan: {
	// Control points to the Layer 2 control being evaluated.
	control: #SingleMapping
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
	requirement: #SingleMapping
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
