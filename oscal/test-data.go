package oscal

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/ossf/gemara"
)

func goodAIGFExample() (gemara.GuidanceDocument, error) {
	testdataPath := "../test-data/good-aigf.yaml"
	data, err := os.ReadFile(testdataPath)
	if err != nil {
		return gemara.GuidanceDocument{}, err
	}
	var l1Docs gemara.GuidanceDocument
	if err := yaml.Unmarshal(data, &l1Docs); err != nil {
		return gemara.GuidanceDocument{}, err
	}
	return l1Docs, nil
}

// guidanceWithExternalExtends returns a guidance document with a guideline that extends an external control.
func guidanceWithExternalExtends() gemara.GuidanceDocument {
	return gemara.GuidanceDocument{
		Title: "Test Guidance",
		Metadata: gemara.Metadata{
			Id:      "TEST-GUIDE",
			Version: "1.0.0",
			Author: gemara.Actor{
				Id:   "test-author",
				Name: "Test Author",
				Type: gemara.Human,
			},
			MappingReferences: []gemara.MappingReference{
				{
					Id:      "NIST-800-53",
					Title:   "NIST SP 800-53",
					Version: "Rev 5",
					Url:     "https://nist.gov/800-53",
				},
			},
		},
		DocumentType: gemara.DocumentType("Framework"),
		Families: []gemara.Family{
			{
				Id:          "AC",
				Title:       "Access Control",
				Description: "Access control family",
			},
		},
		Guidelines: []gemara.Guideline{
			{
				Id:     "TEST-AC-1",
				Title:  "Test Access Control Enhancement",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					ReferenceId: "NIST-800-53",
					EntryId:     "AC-1",
				},
				Objective: "Enhanced access control policy",
				Statements: []gemara.Statement{
					{
						Id:   "1",
						Text: "Additional requirement for access control",
					},
				},
			},
		},
	}
}

// guidanceWithMerging returns a guidance document with multiple guidelines extending the same external control.
func guidanceWithMerging() gemara.GuidanceDocument {
	return gemara.GuidanceDocument{
		Title: "Test Guidance",
		Metadata: gemara.Metadata{
			Id:      "TEST-GUIDE",
			Version: "1.0.0",
			Author: gemara.Actor{
				Id:   "test-author",
				Name: "Test Author",
				Type: gemara.Human,
			},
			MappingReferences: []gemara.MappingReference{
				{
					Id:      "NIST-800-53",
					Title:   "NIST SP 800-53",
					Version: "Rev 5",
					Url:     "https://nist.gov/800-53",
				},
			},
		},
		DocumentType: gemara.DocumentType("Framework"),
		Families: []gemara.Family{
			{
				Id:          "AC",
				Title:       "Access Control",
				Description: "Access control family",
			},
		},
		Guidelines: []gemara.Guideline{
			{
				Id:     "TEST-AC-1",
				Title:  "First Enhancement for AC-1",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					ReferenceId: "NIST-800-53",
					EntryId:     "AC-1",
				},
				Objective: "First objective",
				Statements: []gemara.Statement{
					{
						Id:   "1",
						Text: "First statement",
					},
				},
			},
			{
				Id:     "TEST-AC-1-2",
				Title:  "Second Enhancement for AC-1",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					ReferenceId: "NIST-800-53",
					EntryId:     "AC-1",
				},
				Statements: []gemara.Statement{
					{
						Id:   "2",
						Text: "Second statement",
					},
				},
			},
		},
	}
}

// guidanceWithLocalExtends returns a guidance document with a guideline that extends another guideline in the same document.
func guidanceWithLocalExtends() gemara.GuidanceDocument {
	return gemara.GuidanceDocument{
		Title: "Test Guidance",
		Metadata: gemara.Metadata{
			Id:      "TEST-GUIDE",
			Version: "1.0.0",
			Author: gemara.Actor{
				Id:   "test-author",
				Name: "Test Author",
				Type: gemara.Human,
			},
		},
		DocumentType: gemara.DocumentType("Framework"),
		Families: []gemara.Family{
			{
				Id:          "AC",
				Title:       "Access Control",
				Description: "Access control family",
			},
		},
		Guidelines: []gemara.Guideline{
			{
				Id:     "TEST-AC-1",
				Title:  "Base Access Control",
				Family: "AC",
			},
			{
				Id:     "TEST-AC-1-ENH",
				Title:  "Enhanced Access Control",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					EntryId: "TEST-AC-1",
				},
			},
		},
	}
}

// guidanceWithMultiLevelNested returns a guidance document with multi-level nested local extensions (AC-1 -> AC-1-ENH -> AC-1-ENH-2).
func guidanceWithMultiLevelNested() gemara.GuidanceDocument {
	return gemara.GuidanceDocument{
		Title: "Test Guidance",
		Metadata: gemara.Metadata{
			Id:      "TEST-GUIDE",
			Version: "1.0.0",
			Author: gemara.Actor{
				Id:   "test-author",
				Name: "Test Author",
				Type: gemara.Human,
			},
		},
		DocumentType: gemara.DocumentType("Framework"),
		Families: []gemara.Family{
			{
				Id:    "AC",
				Title: "Access Control",
			},
		},
		Guidelines: []gemara.Guideline{
			{
				Id:       "AC-1",
				Title:    "Base Control",
				Family:   "AC",
				Objective: "Base control objective",
			},
			{
				Id:     "AC-1-ENH",
				Title:  "First Enhancement",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					EntryId: "AC-1",
				},
				Statements: []gemara.Statement{
					{
						Id:   "1",
						Text: "First enhancement statement",
					},
				},
			},
			{
				Id:     "AC-1-ENH-2",
				Title:  "Second Enhancement",
				Family: "AC",
				Extends: &gemara.SingleMapping{
					EntryId: "AC-1-ENH",
				},
				Statements: []gemara.Statement{
					{
						Id:   "1",
						Text: "Second enhancement statement",
					},
				},
			},
		},
	}
}

// guidanceWithImports returns a copy of the provided guidance document with an additional mapping reference added.
func guidanceWithImports(base gemara.GuidanceDocument) gemara.GuidanceDocument {
	mapping := gemara.MappingReference{
		Id:          "EXP",
		Description: "Example",
		Version:     "0.1.0",
		Url:         "https://example.com",
	}
	base.Metadata.MappingReferences = append(base.Metadata.MappingReferences, mapping)
	return base
}
