package schemas

@go(gemara)

#Guidance: {
	metadata?:       #Metadata @go(Metadata)
	title:           string
	"document-type": #DocumentType @go(DocumentType) @yaml("document-type")
	exemptions?: [...#Exemption] @go(Exemptions)

	// Introductory text for the document to be used during rendering
	"front-matter"?: string @go(FrontMatter) @yaml("front-matter,omitempty")
	"categories"?: [...#Category] @go(Categories)
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
	redirect?: #MultiMapping @go(Redirect)
}

// Guideline represents a single guideline within a guidance document
#Guideline: {
	id:         string
	title:      string
	objective?: string
	recommendations?: [...string]
	// Extends allows you to add supplemental guidance within a local guidance document
	// like a control enhancement or from an imported guidance document.
	extends?:   #SingleMapping @go(Extends)
	rationale?: #Rationale     @go(Rationale,optional=nillable)
	statements?: [...#Statement] @go(Statements)
	"guideline-mappings"?: [...#MultiMapping] @go(GuidelineMappings) @yaml("guideline-mappings,omitempty")
	"principle-mappings"?: [...#MultiMapping] @go(PrincipleMappings) @yaml("principle-mappings,omitempty")
	"see-also"?: [...#SingleMapping] @go(SeeAlso) @yaml("see-also,omitempty")
}

// Statement represents a sub-statement within a guideline
#Statement: {
	id:     string
	title?: string
	text:   string
	recommendations?: [...string]
}

// Rationale provides contextual information to help with development and understanding of
// guideline intent.
#Rationale: {
	importance: string
	goals: [...string]
}
