package schemas

@go(gemara)

// Core Document Structure
#PolicyDocument: {
	metadata:          #Metadata
	"organization-id": string @go(OrganizationID) @yaml("organization-id")
	title:             string
	purpose:           string
	contacts:          #Contacts

	scope:                  #Scope
	"implementation-plan"?: #ImplementationPlan @go(ImplementationPlan) @yaml("implementation-plan,omitempty")
	"guidance-references": [...#PolicyMapping] @go(GuidanceReferences) @yaml("guidance-references")
	"control-references": [...#PolicyMapping] @go(ControlReferences) @yaml("control-references")
}

#Contacts: {
	responsible: [...#Contact] // The person or group responsible for implementing controls for technical requirements
	accountable: [...#Contact] // The person or group accountable for evaluating and enforcing the efficacy of technical controls
	consulted?: [...#Contact] // Optional person or group who may be consulted for more information about the technical requirements
	informed?: [...#Contact] // Optional person or group who must recieve updates about compliance with this policy
}

#ImplementationPlan: {
	// The process through which notified parties should be made aware of this policy
	"notification-process"?: string @go(NotificationProcess) @yaml("notification-process",omitempty)
	"notified-parties"?: [...#NotificationGroup] @go(NotifiedParties) @yaml("notified-parties",omitempty)

	evaluation: #ImplementationDetails
	"evaluation-points"?: [...#EvaluationPoint] @go(EvaluationPoints) @yaml("evaluation-points",omitempty)

	enforcement: #ImplementationDetails
	"enforcement-methods"?: [...#EnforcementMethod] @go(EnforcementMethods) @yaml("enforcement-methods",omitempty)

	// The process that will be followed in the event that noncompliance is detected in an applicable resource
	"noncompliance-plan"?: string @go(NoncompliancePlan) @yaml("noncompliance-plan",omitempty)
}

#ImplementationDetails: {
	start: #Datetime
	end?:  #Datetime
	notes: string
}

#Scope: {
	// geopolitical boundaries such as region names or jurisdictions
	boundaries?: [...string]
	// names of technology categories or services
	technologies?: [...string]
	// names of organizations who make the listed technologies available
	providers?: [...string]
}

#PolicyMapping: {
	"reference-id": string @go(ReferenceId) @yaml("reference-id",omitempty)
	"in-scope":     #Scope @go(InScope) @yaml("in-scope",omitempty)
	"out-of-scope": #Scope @go(OutOfScope) @yaml("out-of-scope",omitempty)
	"control-modifications": [...#ControlModifier] @go(ControlModifications) @yaml("control-modifications",omitempty)
	"assessment-requirement-modifications": [...#AssessmentRequirementModifier] @go(AssessmentRequirementModifications) @yaml("assessment-requirement-modifications",omitempty)
	"guideline-modifications": [...#GuidelineModifier] @go(GuidelineModifications) @yaml("guideline-modifications",omitempty)
}

// Modifier Types
#ControlModifier: {
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	title?:     string
	objective?: string
}

#AssessmentRequirementModifier: {
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	text: string
	applicability: [...string]
	recommendation?: string
}

#GuidelineModifier: {
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	title:      string
	objective?: string
	recommendations?: [...string]
	"base-guideline-id"?: string @go(BaseGuidelineID) @yaml("base-guideline-id,omitempty")
	rationale?:           string @go(Rationale,optional=nillable)
	"guideline-mappings"?: [...#Mapping] @go(GuidelineMappings) @yaml("guideline-mappings,omitempty")
	"principle-mappings"?: [...#Mapping] @go(PrincipleMappings) @yaml("principle-mappings,omitempty")
	"see-also"?: [...string] @go(SeeAlso) @yaml("see-also,omitempty")
	"external-references"?: [...string] @go(ExternalReferences) @yaml("external-references,omitempty")
}

#PartModifier: {
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	title?: string
	prose:  string
	recommendations?: [...string]
}

#EvaluationPoint: "development-tools" |
	// For noncompliance risk to workflows or local machines
	"pre-commit-hook" |
	// For noncompliance risk to a development repository
	"pre-merge" |
	// For noncompliance risk to primary repositories
	"pre-build" |
	// For noncompliance risk to built assets
	"pre-release" |
	// For noncompliance risk to released assets
	"pre-deploy" |
	// For noncompliance risk to deployments
	"runtime-adhoc" |
	// For situations where drift may occur
	"runtime-scheduled" |
	// For situations where drift detection is automated
	"runtime-reactive"
// For situations where drift detection is triggered by events

#EnforcementMethod: "Deployment Gate" |
	"Autoremediation" |
	"Manual Remediation"

#NotificationGroup: "Responsible" |
	"Acccountable" |
	"Consulted" |
	"Informed"

#ModType: "increase-strictness" | "clarify" | "reduce-strictness" | "exclude"
