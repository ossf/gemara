package schemas

@go(gemara)

#GuidanceDocument: {
	metadata?:       #Metadata @go(Metadata)
	title:           string
	"document-type": #DocumentType @go(DocumentType) @yaml("document-type")
	exemptions?: [...#Exemption] @go(Exemptions)

	// Introductory text for the document to be used during rendering
	"front-matter"?: string @go(FrontMatter) @yaml("front-matter,omitempty")
	"categories"?: [...#Category] @go(Categories)

	// For inheriting from other guidance documents to create tailored documents/baselines
	"imported-guidelines"?: [...#Mapping] @go(ImportedGuidelines) @yaml("imported-guidelines,omitempty")
	"imported-principles"?: [...#Mapping] @go(ImportedPrinciples) @yaml("imported-principles,omitempty")
}

#DocumentType: "Standard" | "Regulation" | "Best Practice" | "Framework"

// Category represents a logical group of guidelines (i.e. control family)
#Category: {
	id:          string
	title:       string
	description: string
	guidelines?: [...#Guideline]
}

// Exemption represents an exemption with a reason and optional redirect
#Exemption: {
	reason:    string
	redirect?: #Mapping @go(Redirect)
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
