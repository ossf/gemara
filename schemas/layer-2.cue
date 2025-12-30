// Schema lifecycle: experimental | stable | deprecated
@status("experimental")
@if(!stable)
package schemas

@go(gemara)

#Catalog: {
	"metadata"?: #Metadata @go(Metadata)
	title:       string

	families?: [...#Family] @go(Families)
	controls?: [...#Control] @go(Controls)
	threats?: [...#Threat] @go(Threats)
	capabilities?: [...#Capability] @go(Capabilities)

	"imported-controls"?: [...#MultiMapping] @go(ImportedControls)
	"imported-threats"?: [...#MultiMapping] @go(ImportedThreats)
	"imported-capabilities"?: [...#MultiMapping] @go(ImportedCapabilities)
}

#Control: {
	id:        string
	title:     string
	objective: string

	// Family id that this control belongs to
	family: string @go(Family)

	"assessment-requirements": [...#AssessmentRequirement] @go(AssessmentRequirements)
	"guideline-mappings"?: [...#MultiMapping] @go(GuidelineMappings)
	"threat-mappings"?: [...#MultiMapping] @go(ThreatMappings)
}

#Threat: {
	id:          string
	title:       string
	description: string
	capabilities: [...#MultiMapping]

	"external-mappings"?: [...#MultiMapping] @go(ExternalMappings)
}

#Capability: {
	id:          string
	title:       string
	description: string
}

#AssessmentRequirement: {
	id:   string
	text: string
	applicability: [...string]

	recommendation?: string
}
