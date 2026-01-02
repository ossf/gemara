// Schema lifecycle: experimental | stable | deprecated
@status("experimental")
@if(!stable)
package schemas

@go(gemara)

#GuidanceDocument: {
	title:           string
	metadata:        #Metadata     @go(Metadata)
	"document-type": #DocumentType @go(DocumentType) @yaml("document-type")
	// Introductory text for the document to be used during rendering
	"front-matter"?: string @go(FrontMatter) @yaml("front-matter,omitempty")

	families?: [...#Family] @go(Families)
	guidelines?: [...#Guideline] @go(Guidelines)
	exemptions?: [...#Exemption] @go(Exemptions)

	// Guidelines that extend other guidelines must be in the same family as the
	// extended guideline.
	_validateExtensions: {
		for guideline in guidelines if guideline.extends != _|_ {
			if (guideline.extends."reference-id" == "" || guideline.extends."reference-id" == _|_) {
				for extended in guidelines if extended.id == guideline.extends."entry-id" {
					guideline.family == extended.family
				}
			}
		}
	}
}

#DocumentType: "Standard" | "Regulation" | "Best Practice" | "Framework"

// Exemption represents those who are exempt from the full guidance document.
#Exemption: {
	// Description identifies who or what is exempt from the full guidance
	description: string
	// Reason explains why the exemption is granted
	reason: string
	// Redirect points to alternative guidelines or controls that should be followed instead
	redirect?: #MultiMapping @go(Redirect,optional=nillable)
}

// Guideline represents a single guideline within a guidance document
#Guideline: {
	id:         string
	title:      string
	objective?: string

	// Family id that this guideline belongs to
	family: string @go(Family)

	// Maps to fields commonly seen in controls with implementation guidance
	recommendations?: [...string]

	// Extends allows you to add supplemental guidance within a local guidance document
	// like a control enhancement or from an imported guidance document.
	extends?: #SingleMapping @go(Extends,optional=nillable)

	// Applicability specifies the contexts in which this guideline applies.
	applicability?: [...string] @go(Applicability)

	rationale?: #Rationale @go(Rationale,optional=nillable)
	statements?: [...#Statement] @go(Statements)

	"guideline-mappings"?: [...#MultiMapping] @go(GuidelineMappings) @yaml("guideline-mappings,omitempty")
	// A list for associated key principle ids
	"principle-mappings"?: [...#MultiMapping] @go(PrincipleMappings) @yaml("principle-mappings,omitempty")

	// SeeAlso lists related guideline IDs within the same Guidance document.
	"see-also"?: [...string] @go(SeeAlso) @yaml("see-also,omitempty")
}

// Statement represents a structural sub-requirement within a guideline
// They do not increase strictness and all statements within a guideline apply together.
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
