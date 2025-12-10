package schemas

@go(gemara)

// Core Document Structure
#Policy: {
	metadata:          #Metadata
	"organization-id": string @go(OrganizationID) @yaml("organization-id")
	title:             string
	purpose:           string
	contacts:          #Contacts

	scope: #Scope
	// ImportedPolicies: References to other Layer 3 Policy documents for policy composition.
	// When a policy imports another policy via ImportedPolicies:
	//   * ImplementationPlan: Inherited, evaluators are ADDITIVE (union) if child also defines evaluators
	// - Merged reference sets (can be modified):
	//   * ControlReferences: Merged from imported policy, can be modified via ControlModifications
	//   * GuidanceReferences: Merged from imported policy, can be modified via GuidelineModifications
	// - Modifications chain sequentially (e.g., Base → Parent → Child), preserving all modifications
	"imported-policies"?: [...#PolicyMapping] @go(ImportedPolicies) @yaml("imported-policies,omitempty")
	"implementation-plan"?: #ImplementationPlan @go(ImplementationPlan) @yaml("implementation-plan,omitempty")
	// GuidanceReferences: References to Layer 1 Guidance documents (abstract, high-level guidance).
	// Modifications chain sequentially (e.g., Guidance → Parent → Child), preserving all modifications.
	"guidance-references"?: [...#PolicyMapping] @go(GuidanceReferences) @yaml("guidance-references")
	// ControlReferences: References to Layer 2 Control catalogs (technology-specific, threat-informed controls).
	// Modifications chain sequentially (e.g., Catalog → Parent → Child), preserving all modifications.
	"control-references"?: [...#PolicyMapping] @go(ControlReferences) @yaml("control-references")
}

#Contacts: {
	responsible: [...#Contact] // The person or group responsible for implementing controls for technical requirements
	accountable: [...#Contact] // The person or group accountable for evaluating and enforcing the efficacy of technical controls
	consulted?: [...#Contact] // Optional person or group who may be consulted for more information about the technical requirements
	informed?: [...#Contact] // Optional person or group who must receive updates about compliance with this policy
}

#ImplementationPlan: {
	// The process through which notified parties should be made aware of this policy
	"notification-process"?: string @go(NotificationProcess) @yaml("notification-process,omitempty")
	"notified-parties"?: [...#NotificationGroup] @go(NotifiedParties) @yaml("notified-parties,omitempty")

	"evaluation-timeline": #ImplementationDetails @go(EvaluationTimeline) @yaml("evaluation-timeline")
	// Evaluators: Actors (human or software) that perform assessments.
	// When importing policies: Inherited from imported policy, evaluators are ADDITIVE (union) if child also defines evaluators.
	// This prevents broken evaluator references in assessment requirements.
	evaluators?: [...#Actor] @go(Evaluators) @yaml("evaluators,omitempty")

	"enforcement-timeline": #ImplementationDetails @go(EnforcementTimeline) @yaml("enforcement-timeline")
	"enforcement-methods"?: [...#EnforcementMethod] @go(EnforcementMethods) @yaml("enforcement-methods,omitempty")

	// The consequence that will be applied in the event that noncompliance is detected
	"noncompliance-consequence"?: string @go(NoncomplianceConsequence) @yaml("noncompliance-consequence,omitempty")
}

#ImplementationDetails: {
	start:  #Datetime
	end?:   #Datetime
	notes?: string
}

// Scope is descriptive metadata for tools.
// When importing policies: Inherited from imported policy, but child's scope OVERRIDES parent's if defined.
#Scope: {
	// geopolitical boundaries such as region names or jurisdictions
	boundaries?: [...string]
	// names of technology categories or services
	technologies?: [...string]
	// names of organizations who make the listed technologies available
	providers?: [...string]
}

// Layer 3 specific mapping that extends common Mapping with modifications.
// Used by ImportedPolicies, GuidanceReferences, and ControlReferences.
// Modifications chain sequentially when importing policies (e.g., Base → Parent → Child), preserving all modifications.
#PolicyMapping: {
	"reference-id": string @go(ReferenceId) @yaml("reference-id")
	// ControlModifications: Modify controls from referenced catalog/policy.
	// When used in ImportedPolicies: Modifies controls from imported policy's ControlReferences (sequential chain).
	// When used in ControlReferences: Modifies controls from referenced catalog (sequential chain).
	"control-modifications"?: [...#ControlModifier] @go(ControlModifications) @yaml("control-modifications,omitempty")
	// AssessmentRequirementModifications: Modify assessment requirements (which have evaluators).
	// When used in ImportedPolicies: Modifies assessment requirements from imported policy's ControlReferences (sequential chain).
	// When used in ControlReferences: Modifies assessment requirements from referenced catalog (sequential chain).
	"assessment-requirement-modifications"?: [...#AssessmentRequirementModifier] @go(AssessmentRequirementModifications) @yaml("assessment-requirement-modifications,omitempty")
	// GuidelineModifications: Modify guidelines from referenced guidance document.
	// When used in ImportedPolicies: Modifies guidelines from imported policy's GuidanceReferences (sequential chain).
	// When used in GuidanceReferences: Modifies guidelines from referenced guidance document (sequential chain).
	"guideline-modifications"?: [...#GuidelineModifier] @go(GuidelineModifications) @yaml("guideline-modifications,omitempty")
	// StatementModifications: Modify statements within guidelines.
	// When used in ImportedPolicies: Modifies statements from imported policy's GuidanceReferences (sequential chain).
	// When used in GuidanceReferences: Modifies statements from referenced guidance document (sequential chain).
	"statement-modifications"?: [...#StatementModifier] @go(StatementModifications) @yaml("statement-modifications,omitempty")
}

// Modifier Types
#ControlModifier: {
	id?:                      string   @go(Id)
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	overrides?:  #Control           @go(Overrides,optional=nillable)
	extensions?: #ControlExtensions @go(Extensions,optional=nillable)
}

#ControlExtensions: {
	severity?:                   #Severity @go(Severity)
	"auto-remediation-allowed"?: bool      @go(AutoRemediationAllowed) @yaml("auto-remediation-allowed,omitempty")
	"deployment-gate-allowed"?:  bool      @go(DeploymentGateAllowed) @yaml("deployment-gate-allowed,omitempty")
}

// Severity represents the severity level of a control
#Severity: "Critical" | "High" | "Medium" | "Low" | "Info" | "Unknown" @go(-)

#AssessmentRequirementModifier: {
	id?:                      string   @go(Id)
	"target-id":              string   @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType @go(ModType) @yaml("modification-type")
	"modification-rationale": string   @go(ModificationRationale) @yaml("modification-rationale")

	overrides?:  #AssessmentRequirement           @go(Overrides,optional=nillable)
	extensions?: #AssessmentRequirementExtensions @go(Extensions,optional=nillable)
}

#AssessmentRequirementExtensions: {
	"required-evaluators"?: [...string] @go(RequiredEvaluators) @yaml("required-evaluators,omitempty")
	"optional-evaluators"?: [...string] @go(OptionalEvaluators) @yaml("optional-evaluators,omitempty")
	"evidence-requirements"?: string              @go(EvidenceRequirements) @yaml("evidence-requirements,omitempty")
	"resolution-strategy"?:   #ResolutionStrategy @go(ResolutionStrategy) @yaml("resolution-strategy,omitempty")
	"evaluation-points"?: [...#EvaluationPoint] @go(EvaluationPoints) @yaml("evaluation-points,omitempty")
}

#GuidelineModifier: {
	id?:                      string     @go(Id)
	"target-id":              string     @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType   @go(ModType) @yaml("modification-type")
	"modification-rationale": string     @go(ModificationRationale) @yaml("modification-rationale")
	overrides?:               #Guideline @go(Overrides,optional=nillable)
}

#StatementModifier: {
	id?:                      string     @go(Id)
	"target-id":              string     @go(TargetId) @yaml("target-id")
	"modification-type":      #ModType   @go(ModType) @yaml("modification-type")
	"modification-rationale": string     @go(ModificationRationale) @yaml("modification-rationale")
	overrides?:               #Statement @go(Overrides,optional=nillable)
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
	"Accountable" |
	"Consulted" |
	"Informed"

// ModType: Semantic modification types for policy tailoring.
#ModType: "increase-strictness" | "clarify" | "reduce-strictness" | "exclude"

// ResolutionStrategy defines how to resolve conflicts when multiple evaluators produce different results
// for the same assessment requirement. Options:
// - "MostSevere": Use the most severe result from all evaluators (Failed > Unknown > NeedsReview > Passed)
// - "ManualOverride": Give precedence to manual review evaluators over automated evaluators when results conflict
// - "AuthoritativeConfirmation": Require confirmation from authoritative evaluators before triggering findings from non-authoritative evaluators
#ResolutionStrategy: "MostSevere" | "ManualOverride" | "AuthoritativeConfirmation" @go(-)
