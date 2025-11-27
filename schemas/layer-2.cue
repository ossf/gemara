package schemas

@go(gemara)

#Catalog: {
	"metadata"?: #Metadata @go(Metadata)
	title:       string

	"control-families"?: [...#ControlFamily] @go(ControlFamilies)
	threats?: [...#Threat] @go(Threats)
	capabilities?: [...#Capability] @go(Capabilities)

	"imported-controls"?: [...#Mapping] @go(ImportedControls)
	"imported-threats"?: [...#Mapping] @go(ImportedThreats)
	"imported-capabilities"?: [...#Mapping] @go(ImportedCapabilities)
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
	"guideline-mappings"?: [...#Mapping] @go(GuidelineMappings)
	"threat-mappings"?: [...#Mapping] @go(ThreatMappings)
}

#Threat: {
	id:          string
	title:       string
	description: string
	capabilities: [...#Mapping]

	"external-mappings"?: [...#Mapping] @go(ExternalMappings)
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
