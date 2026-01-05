// Schema lifecycle: experimental | stable | deprecated
@status("stable")

package schemas

// ============================================================================
// Mapping Types - MappingReference, MappingEntry, MultiMapping, SingleMapping
// ============================================================================

// MappingReference represents a reference to an external document with full metadata.
#MappingReference: {
	id:           string
	title:        string
	version:      string
	description?: string
	url?:         =~"^(https?|file)://[^\\s]+$"
}

// MultiMapping represents a mapping to an external reference with one or more entries.
#MultiMapping: {
	// ReferenceId should reference the corresponding MappingReference id from metadata
	"reference-id": string @go(ReferenceId)
	entries: [#MappingEntry, ...#MappingEntry] @go(Entries)
	remarks?: string
}

// SingleMapping represents how a specific entry (control/requirement/procedure) maps to a MappingReference.
#SingleMapping: {
	// ReferenceId should reference the corresponding MappingReference id from metadata
	"reference-id"?: string @go(ReferenceId)
	"entry-id":      string @go(EntryId)
	remarks?:        string
}

// MappingEntry represents a single entry within a mapping
#MappingEntry: {
	"reference-id": string @go(ReferenceId)
	// Strength quantifies the degree of correlation or relationship between the mapped items.
	// Range: 1-10. Zero value means not yet quantified.
	strength?: int & >=1 & <=10
	remarks?:  string
}
