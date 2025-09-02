package schemas

import "time"

#EvaluationResults: {
	"evaluation-set": [#ControlEvaluation, ...#ControlEvaluation] @go(EvaluationSet)
	...
}

#EvaluationPlan: {
	metadata: #Metadata
	plans: [...#AssessmentPlan]
}

// AssessmentPlan defines all testing procedures for a control id.
#AssessmentPlan: {
	"control-id": string @go(ControlId)
	"assessments": [...#Assessment] @go(Assessments)
}

#ControlEvaluation: {
	name:              string
	"control-id":      string @go(ControlId)
	result:            #Result
	message:           string
	"corrupted-state": bool @go(CorruptedState)
	"assessment-logs": [...#AssessmentLog] @go(AssessmentLogs)
}

// AssessmentPlan defines all testing procedures for a requirement.
#Assessment: {
	// RequirementID is the unique identifier for the requirement being tested
	"requirement-id": string @go(RequirementId)
	// Procedures defined possible testing procedures to evaluate the requirement.
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

#AssessmentLog: {
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

#Metadata: {
	id:       string
	version?: string
	author:   #Contact
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
