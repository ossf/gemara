package schemas

import "time"

// #EnforcementAction defines an auditable record of policy enforcement.
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
	// Control defines the Layer 2 control being enforced.
	control: #Mapping
	// Findings defines Layer 4 AssessmentLog outcomes that triggered this enforcement action.
	findings?: [...#Finding]
	// Exception defines an optional exception that may apply to this finding.
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

// Metadata contains metadata about the Layer 5 Enforcement Action.
#Metadata: {
	// Id defines the unique identifier for the metadata record.
	id: string
	// Version defines the version of the metadata schema.
	version?: string
	// Author defines the entity that produced the enforcement action.
	author: #Author
	// MappingReferences defines references to external standards, frameworks, or documents.
	"mapping-references"?: [...#MappingReference] @go(MappingReferences) @yaml("mapping-references,omitempty")
}

// Contact represents contact information for a person.
#Contact: {
	// Name defines the contact person's name.
	name: string
	// Primary indicates whether this contact is the first point of contact for inquiries. Only one entry should be marked as primary.
	primary: bool
	// Affiliation defines the entity with which the contact is affiliated, such as a school or employer.
	affiliation?: string @go(Affiliation,type=*string)
	// Email defines a preferred email address to reach the contact.
	email?: #Email @go(Email,type=*Email)
	// Social defines a social media handle or profile for the contact.
	social?: string @go(Social,type=*string)
}

// Author contains the information about the entity that produced the evaluation plan or log.
#Author: {
	// Name defines the name of the author.
	name: string
	// Uri defines a URI for the author.
	uri?: string
	// Version defines the version of the authoring entity.
	version?: string
	// Contact defines contact information for the author.
	contact?: #Contact @go(Contact)
}

// MappingReference provides references to external standards, frameworks, or documents.
#MappingReference: {
	// Id defines the unique identifier for the referenced document or standard.
	id: string
	// Title defines a human-readable title of the referenced document.
	title: string
	// Version defines the version of the referenced document.
	version: string
	// Description defines an optional description of the referenced document.
	description?: string
	// Url defines an optional URL to access the referenced document (must be valid HTTP/HTTPS URL).
	url?: =~"^https?://[^\\s]+$"
}

// Mapping represents a mapping between internal controls and external standards.
#Mapping: {
	// ReferenceId defines the corresponding MappingReference id.
	"reference-id": string @go(ReferenceId)
	// EntryId defines the specific element within the referenced document.
	"entry-id": string @go(EntryId)
	// Strength defines how effectively the referenced item addresses the associated control or procedure on a scale of 1 to 10, with 10 being the most effective.
	strength?: int & >=1 & <=10
	// Remarks defines additional context about the mapping entry.
	remarks?: string
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
	"compensating-controls"?: [...#Mapping] @go(CompensatingControls)
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

// Finding is a self-contained record of a detected issue.
#Finding: {
	// Requirement defines the specific requirement that was evaluated.
	"requirement": #Mapping
	// Result defines the result of evaluating this requirement.
	result: #Result
	// Message defines a human-readable description of what was found.
	message: string
}

// Datetime represents a timestamp in ISO 8601 format with timezone.
#Datetime: time.Format("2006-01-02T15:04:05Z07:00") @go(Datetime,format="date-time")

// Email represents a valid email address format.
#Email: =~"^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Za-z]{2,}$"

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
