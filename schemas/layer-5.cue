package schemas

@go(gemara)

// #EnforcementAction defines an auditable record of enforcement actions taken to ensure
// compliance with Layer 3 policy requirements based on Layer 4 evaluation findings.
// Layer 5 links to Layer 3 (Policy) and Layer 4 (Findings via requirement/procedure mappings).
// Layer 2 (Controls) are accessible through findings' requirement mappings.
#EnforcementAction: {
	metadata: #Metadata
	// Executed indicates whether the enforcement action was successfully executed.
	executed: bool
	// ExecutedAt defines when the enforcement action was executed.
	"executed-at": #Datetime @go(ExecutedAt)
	// Message defines a brief description of what enforcement action was actually taken.
	message?: string
	// Target defines the subject of the enforcement action.
	target?: #Target
	// Action defines the high-level action performed during enforcement.
	action: #Action
	// Policy defines the Layer 3 policy document whose requirements this enforcement action
	// ensures compliance with.
	policy: #SingleMapping
	// Findings defines Layer 4 AssessmentLog outcomes that triggered this enforcement action.
	// Findings are required as they inform the action taken.
	// Findings reference Layer 4 data through their requirement and procedure mappings.
	// The Layer 2 control being enforced can be identified through the findings' requirement mappings.
	findings: [...#Finding]
	// Exception defines an optional exception that applies to all findings in this enforcement action.
	// When different exceptions are needed for different findings, create separate EnforcementAction records.
	exception?: #Exception
	// RemediationPlan uniquely identifies the remediation response when the Layer 3 Enforcement
	// Method is AutoRemediation.
	"remediation-plan"?: string @go(RemediationPlanId)
	// NotificationPlan uniquely identifies the notification response when the Layer 3 Enforcement
	// Method is Manual Remediation.
	"notification-plan"?: string @go(NotificationPlan)
	// EnforcementPlan uniquely identifies the enforcement response when the Layer 3 Enforcement
	// Method is Deployment Gate.
	"enforcement-plan"?: string @go(EnforcementPlan)
}

// Exception represents an approved exception to policy enforcement.
#Exception: {
	// Id defines the unique identifier for this exception.
	id: string
	// ApprovedBy defines the person or entity who approved this exception.
	"approved-by": #Contact @go(ApprovedBy)
	// ApprovalDate defines the date and time when the exception was approved.
	"approval-date": #Datetime @go(ApprovalDate)
	// ExpirationDate defines the optional date when this exception expires.
	"expiration-date"?: #Datetime @go(ExpirationDate)
	// Justification defines the justification for why this exception is necessary.
	"justification": string @go(Justification)
	// RiskLevel defines the risk level associated with this exception.
	"risk-level": #RiskLevel @go(RiskLevel) @yaml("risk-level")
	// CompensatingControls defines an optional list of compensating controls implemented to mitigate risk.
	"compensating-controls"?: [...#MultiMapping] @go(CompensatingControls)
	// ReviewDate defines the optional date when this exception should be reviewed.
	"review-date"?: #Datetime @go(ReviewDate)
}

// Target defines the subject of the enforcement action.
#Target: {
	// TargetId uniquely identifies the specific target instance.
	"target-id": string @go(TargetId)
	// TargetName defines a human-readable name of the target.
	"target-name": string @go(TargetName)
	// TargetType defines the type or category of the target.
	"target-type": string @go(TargetType)
	// Environment defines the environment where the target exists.
	environment?: string
}

// Finding represents Layer 5's opinion about policy conformance based on Layer 4 evidence.
// Findings reference Layer 4 AssessmentLogs (evidence) and form an opinion about that evidence.
#Finding: {
	// Logs references the Layer 4 AssessmentLog entries that serve as evidence for this finding.
	// These logs contain the actual evidence of policy conformance.
	logs: [...#AssessmentLogReference]
	// Result defines Layer 5's opinion about the evidence from the referenced logs.
	// This is the result of evaluating the evidence using conflict resolution strategies.
	result: #Result
	// Message defines a human-readable description of Layer 5's opinion about what was found.
	message: string
	// WeightedScore defines the calculated weighted score from evaluator results.
	// This score is computed using CVSS-inspired weighted averaging, incorporating
	// strategy-based weights and confidence levels.
	// Lower scores indicate more severe findings (0.0 = Failed, 3.0+ = Passed).
	"weighted-score"?: float @go(WeightedScore)
}

// AssessmentLogReference identifies a Layer 4 AssessmentLog that serves as evidence.
// This reference allows Layer 5 findings to link to the evidence rather than copy it.
#AssessmentLogReference: {
	// EvaluatorId identifies the evaluator that produced the assessment log.
	"evaluator-id": string @go(EvaluatorId)
	// Requirement identifies the requirement that was evaluated in the log.
	// This matches AssessmentLog.Requirement.
	requirement: #SingleMapping
	// StartTime identifies when the assessment began, used to uniquely identify the log instance.
	// This matches AssessmentLog.Start.
	"start-time": #Datetime @go(StartTime)
	// LogId optionally provides a unique identifier for the log if one exists in the system.
	// This can be used when logs are stored in a system that assigns unique IDs.
	"log-id"?: string @go(LogId)
}

// Action is the high-level enforcement outcome.
#Action:
	"Block" |
	"Allow" |
	"Remediate" |
	"Waive" |
	"Notify" |
	"Unknown"

// Result is the outcome of the assessment.
#Result: "Not Run" | "Passed" | "Failed" | "Needs Review" | "Not Applicable" | "Unknown"

// RiskLevel from Layer 3 (Policy layer)
#RiskLevel: "Critical" | "High" | "Medium" | "Low" | "Informational"
