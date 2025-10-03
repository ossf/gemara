package schemas

@go(layer1)

#GuidanceDocument: {
	metadata?: #Metadata

	// Introductory text for the document to be used during rendering
	"front-matter"?: string @go(FrontMatter) @yaml("front-matter,omitempty")
	"categories"?: [...#Category] @go(Categories)

	// For inheriting from other guidance documents to create tailored documents/baselines
	"imported-guidelines"?: [...#Mapping] @go(ImportedGuidelines) @yaml("imported-guidelines,omitempty")
	"imported-principles"?: [...#Mapping] @go(ImportedPrinciples) @yaml("imported-principles,omitempty")
}

#Metadata: {
	id:                  string
	title:               string
	description:         string
	author:              string
	version?:            string
	"last-modified"?:    string @go(LastModified) @yaml("last-modified,omitempty")
	"publication-date"?: string @go(PublicationDate) @yaml("publication-date,omitempty")

	"mapping-references"?: [...#MappingReference] @go(MappingReferences) @yaml("mapping-references,omitempty")

	"document-type"?: #DocumentType  @go(DocumentType)
	applicability?:   #Applicability @go(Applicability,optional=nillable)
	exemptions?: [...string]
}

#DocumentType: "Standard" | "Regulation" | "Best Practice" | "Framework"

#Applicability: {
	// Inclusion by geographical or legal areas
	jurisdictions?: [...string]
	// Inclusion by types of technology or technological environments
	"technology-domains"?: [...string] @go(TechnologyDomains) @yaml("technology-domains,omitempty")
	// Inclusion by industry sectors or verticals
	"industry-sectors"?: [...string] @go(IndustrySectors) @yaml("industry-sectors,omitempty")
}

// Category represents a logical group of guidelines (i.e. control family)
#Category: {
	id:          string
	title:       string
	description: string
	guidelines?: [...#Guideline]
}

// Rationale provides contextual information to help with development and understanding of
// guideline intent.
#Rationale: {
	// Negative results expected from the guideline's lack of implementation
	risks: [...#Risk]
	// Positive results expected from the guideline's implementation
	outcomes: [...#Outcome]
}

#Risk: {
	title:       string
	description: string
}

#Outcome: {
	title:       string
	description: string
}

#Guideline: {
	id:         string
	title:      string
	objective?: string

	// Maps to fields commonly seen in controls with implementation guidance
	recommendations?: [...string]

	// For control enhancements (ex. AC-2(1) in 800-53)
	// The base-guideline-id is needed to achieve full context for the enhancement
	"base-guideline-id"?: string @go(BaseGuidelineID) @yaml("base-guideline-id,omitempty")

	rationale?: #Rationale @go(Rationale,optional=nillable)

	// Represents individual guideline parts/statements
	"guideline-parts"?: [...#Part] @go(GuidelineParts) @yaml("guideline-parts,omitempty")
	// Crosswalking this guideline to other guidelines in other documents
	"guideline-mappings"?: [...#Mapping] @go(GuidelineMappings) @yaml("guideline-mappings,omitempty")
	// A list for associated key principle ids
	"principle-mappings"?: [...#Mapping] @go(PrincipleMappings) @yaml("principle-mappings,omitempty")

	// This is akin to related controls, but using more explicit terminology
	"see-also"?: [...string] @go(SeeAlso) @yaml("see-also,omitempty")
}

// Parts include sub-statements of a guideline that can be assessed individually
#Part: {
	id:     string
	title?: string
	text:   string
	recommendations?: [...string]
}

// Mapping references is the same from Layer2, but intended for Layer 1 to Layer 1 mappings
// instead of Layer 2 to Layer 1 mappings.
#MappingReference: {
	id:           string
	title:        string
	version:      string
	description?: string
	issuer?:      string
	url?:         =~"^https?://[^\\s]+$"
}

#Mapping: {
	"reference-id": string @go(ReferenceId)
	entries?: [...#MappingEntry] @go(Entries,optional=nillable)
	// Adding context about this particular mapping and why it was mapped.
	remarks?: string
}

#MappingEntry: {
	"reference-id": string @go(ReferenceId)
	strength:       int & >=1 & <=10
	remarks?:       string
}
