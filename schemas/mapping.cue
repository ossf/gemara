package schemas

// ============================================================================
// Mapping Types - MappingReference, MappingEntry, Mapping, EntryMapping
// ============================================================================

// MappingReference represents a reference to an external document with full metadata.
#MappingReference: {
	id:           string
	title:        string
	version:      string
	description?: string
	url?:         =~"^(https?|file)://[^\\s]+$"
}

// Mapping represents a mapping to an external reference with one or more entries.
#Mapping: {
	// ReferenceId should reference the corresponding MappingReference id from metadata
	"reference-id": string @go(ReferenceId)
	entries: [#MappingEntry, ...#MappingEntry] @go(Entries)
	remarks?: string
}

// EntryMapping represents how a specific entry (control/requirement/procedure) maps to a MappingReference.
#EntryMapping: {
	// ReferenceId should reference the corresponding MappingReference id from metadata
	"reference-id"?: string @go(ReferenceId)
	"entry-id":      string @go(EntryId)
	remarks?:        string
}

// MappingEntry represents a single entry within a mapping
#MappingEntry: {
	"reference-id": string @go(ReferenceId)
	strength:       int & >=1 & <=10
	remarks?:       string
}
