// Schema lifecycle: experimental | stable | deprecated
@status("stable")
package schemas

// Metadata represents common metadata fields shared across all layers
#Metadata: {
	id:          string
	version?:    string
	date?:       #Date @go(Date)
	description: string
	author:      #Actor
	"mapping-references"?: [...#MappingReference] @go(MappingReferences) @yaml("mapping-references,omitempty")
	"applicability-categories"?: [...#Category] @go(ApplicabilityCategories) @yaml("applicability-categories,omitempty")
	draft?:   bool
	lexicon?: string
}
