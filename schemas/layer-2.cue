package schemas

@go(gemara)
// @status tracks schema lifecycle: experimental | stable | deprecated
@status("experimental")

#Catalog: {
	"metadata"?: #Metadata @go(Metadata)
	title:       string

	"control-families"?: [...#ControlFamily] @go(ControlFamilies)
	threats?: [...#Threat] @go(Threats)
	capabilities?: [...#Capability] @go(Capabilities)

	"imported-controls"?: [...#MultiMapping] @go(ImportedControls)
	"imported-threats"?: [...#MultiMapping] @go(ImportedThreats)
	"imported-capabilities"?: [...#MultiMapping] @go(ImportedCapabilities)
}

#ControlFamily: {
	id:          string
	title:       string
	description: string
	controls: [...#Control]
}

#Control: {
	id:        string
	title:     string
	objective: string
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
