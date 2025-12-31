package oscal

import (
	"os"
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/ossf/gemara"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

func TestFromGuidance(t *testing.T) {
	goodAIFG, err := goodAIGFExample()
	require.NoError(t, err)

	tests := []struct {
		name        string
		guidance    gemara.GuidanceDocument
		catalogHref string
		wantErr     bool
	}{
		{
			name:        "valid guidance document",
			guidance:    goodAIFG,
			catalogHref: "test-catalog.json",
			wantErr:     false,
		},
		{
			name:        "error when catalogHref is empty",
			guidance:    goodAIFG,
			catalogHref: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog, profile, err := FromGuidance(&tt.guidance, tt.catalogHref)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Validate catalog
			catalogModel := oscalTypes.OscalModels{Catalog: &catalog}
			assert.NoError(t, oscalUtils.Validate(catalogModel))
			assert.NotEmpty(t, catalog.UUID)

			// Validate profile
			profileModel := oscalTypes.OscalModels{Profile: &profile}
			assert.NoError(t, oscalUtils.Validate(profileModel))
			assert.NotEmpty(t, profile.UUID)

			// Verify catalogHref is in imports
			assert.NotEmpty(t, profile.Imports)
			hasProvidedImport := false
			for _, imp := range profile.Imports {
				if imp.Href == tt.catalogHref && imp.IncludeAll != nil {
					hasProvidedImport = true
					break
				}
			}
			assert.True(t, hasProvidedImport, "profile should have provided catalogHref in imports")
		})
	}
}

func TestToCatalogFromGuidance(t *testing.T) {
	goodAIFG, err := goodAIGFExample()
	require.NoError(t, err)

	tests := []struct {
		name       string
		guidance   gemara.GuidanceDocument
		wantGroups []oscalTypes.Group
		wantErr    bool
	}{
		{
			name:     "Good AIGF",
			guidance: goodAIFG,
			wantGroups: []oscalTypes.Group{
				{
					Class: "family",
					ID:    "DET",
					Title: "Detective",
					Controls: &[]oscalTypes.Control{
						{
							Class: "FINOS-AIR",
							ID:    "air-det-011",
							Title: "Human Feedback Loop for AI Systems",
							Links: &[]oscalTypes.Link{
								{
									Href: "#air-det-015",
									Rel:  "related",
								},
								{
									Href: "#air-det-004",
									Rel:  "related",
								},
								{
									Href: "#air-prev-005",
									Rel:  "related",
								},
								{
									Href: "#placeholder",
									Rel:  "reference",
								},
								{
									Href: "#placeholder",
									Rel:  "reference",
								},
							},
							Parts: &[]oscalTypes.Part{
								{
									Name: "statement",
									ID:   "air-det-011_smt",
									Parts: &[]oscalTypes.Part{
										{
											Name:  "item",
											ID:    "air-det-011_smt.1",
											Title: "Designing the Feedback Mechanism",
											Prose: "Implementing an effective human feedback loop involves careful design of the mechanism.",
										},
										{
											Name:  "item",
											ID:    "air-det-011_smt.2",
											Title: "Types of Feedback and Collection Methods",
											Prose: "Implementing an effective human feedback loop involves clear collection processes.",
										},
									},
								},
								{
									Name: "assessment-objective",
									ID:   "air-det-011_obj",
									Parts: &[]oscalTypes.Part{
										{
											Name: "assessment-objective",
											ID:   "air-det-011_obj.1",
											Links: &[]oscalTypes.Link{
												{
													Href: "#air-det-011_smt.1",
													Rel:  "assessment-for",
												},
											},
											Prose: "Define Intended Use and KPIs:\nObjectives: Clearly document how feedback data will be utilized, such as for prompt fine-tuning, RAG document updates,model/data drift detection, " +
												"or more advanced uses like Reinforcement Learning from Human Feedback (RLHF).\nKPI Alignment: Design feedback questions and metrics to align with the solution's key performance indicators " +
												"(KPIs). For example, if accuracy is a KPI, feedback might involve users or SMEs annotating if an answer was correct.",
										},
										{
											Name: "assessment-objective",
											ID:   "air-det-011_obj.2",
											Links: &[]oscalTypes.Link{
												{
													Href: "#air-det-011_smt.2",
													Rel:  "assessment-for",
												},
											},
											Prose: "Quantitative Feedback:\nDescription: Involves collecting structured responses that can be easily aggregated and measured, such as numerical ratings (e.g., \"Rate this response on " +
												"a scale of 1-5 for helpfulness\"), categorical choices (e.g., \"Was this answer: Correct/Incorrect/Partially Correct\"), or binary responses (e.g., thumbs up/down)." +
												"\nUse Cases: Effective for tracking trends, measuring against KPIs, and quickly identifying areas of high or low performance.",
										},
									},
								},
								{
									Name: "overview",
									ID:   "air-det-011_ovw",
									Prose: "A Human Feedback Loop is a critical detective and continuous improvement mechanism that involves systematically collecting, analyzing, and acting upon feedback provided by human users, " +
										"subject matter experts (SMEs), or reviewers regarding an AI system's performance, outputs, or behavior.",
								},
							},
						},
						{
							Class: "FINOS-AIR",
							ID:    "air-det-004",
							Title: "Example Detective Control 004",
							Parts: &[]oscalTypes.Part{
								{
									Name: "statement",
									ID:   "air-det-004_smt",
								},
								{
									Name: "assessment-objective",
									ID:   "air-det-004_obj",
								},
								{
									Name:  "overview",
									ID:    "air-det-004_ovw",
									Prose: "Placeholder control for testing references.",
								},
							},
						},
						{
							Class: "FINOS-AIR",
							ID:    "air-det-015",
							Title: "Example Detective Control 015",
							Parts: &[]oscalTypes.Part{
								{
									Name: "statement",
									ID:   "air-det-015_smt",
								},
								{
									Name: "assessment-objective",
									ID:   "air-det-015_obj",
								},
								{
									Name:  "overview",
									ID:    "air-det-015_ovw",
									Prose: "Placeholder control for testing references.",
								},
							},
						},
					},
				},
				{
					Class: "family",
					ID:    "PREV",
					Title: "Preventive",
					Controls: &[]oscalTypes.Control{
						{
							Class: "FINOS-AIR",
							ID:    "air-prev-005",
							Title: "Example Preventive Control 005",
							Parts: &[]oscalTypes.Part{
								{
									Name: "statement",
									ID:   "air-prev-005_smt",
								},
								{
									Name: "assessment-objective",
									ID:   "air-prev-005_obj",
								},
								{
									Name:  "overview",
									ID:    "air-prev-005_ovw",
									Prose: "Placeholder control for testing references.",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Failure/EmptyGuidance",
			guidance: gemara.GuidanceDocument{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog, _, err := FromGuidance(&tt.guidance, "test-catalog.json")
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				catalogModel := oscalTypes.OscalModels{Catalog: &catalog}
				err = oscalUtils.Validate(catalogModel)
				assert.NoError(t, err)
				// Sort slices to ignore order when comparing
				sortGroups := cmpopts.SortSlices(func(a, b oscalTypes.Group) bool {
					return a.ID < b.ID
				})
				sortControls := cmpopts.SortSlices(func(a, b oscalTypes.Control) bool {
					return a.ID < b.ID
				})
				if diff := cmp.Diff(tt.wantGroups, *catalog.Groups, cmpopts.IgnoreFields(oscalTypes.Link{}, "Href"), sortGroups, sortControls); diff != "" {
					t.Errorf("group mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestToProfileFromGuidance(t *testing.T) {
	goodAIFG, err := goodAIGFExample()
	require.NoError(t, err)

	guidanceWithImports := goodAIFG
	// Add some shared guidelines
	mapping := gemara.MappingReference{
		Id:          "EXP",
		Description: "Example",
		Version:     "0.1.0",
		Url:         "https://example.com",
	}

	importedGuidelines := gemara.MultiMapping{
		ReferenceId: "EXP",
		Entries: []gemara.MappingEntry{
			{
				ReferenceId: "EX-1",
			},
			{
				// Intentionally adding a control that
				// needs to be normalized
				ReferenceId: "EX-1(2)",
			},
			{
				ReferenceId: "EX-2",
			},
		},
	}
	guidanceWithImports.Metadata.MappingReferences = append(guidanceWithImports.Metadata.MappingReferences, mapping)
	guidanceWithImports.ImportedGuidelines = append(guidanceWithImports.ImportedGuidelines, importedGuidelines)

	tests := []struct {
		name        string
		guidance    gemara.GuidanceDocument
		options     []GenerateOption
		wantImports []oscalTypes.Import
	}{
		{
			name:     "Success/LocalOnly",
			guidance: goodAIFG,
			wantImports: []oscalTypes.Import{
				{
					Href:       "testHref",
					IncludeAll: &oscalTypes.IncludeAll{},
				},
			},
		},
		{
			name:     "Success/WithImports",
			guidance: guidanceWithImports,
			wantImports: []oscalTypes.Import{
				{
					Href: "https://example.com",
					IncludeControls: &[]oscalTypes.SelectControlById{
						{
							WithIds: &[]string{
								"ex-1",
								"ex-1.2",
								"ex-2",
							},
						},
					},
				},
				{
					Href:       "testHref",
					IncludeAll: &oscalTypes.IncludeAll{},
				},
			},
		},
		{
			name:     "Success/WithImportOverride",
			guidance: guidanceWithImports,
			options: []GenerateOption{
				WithOSCALImports(map[string]string{
					"EXP": "https://example.com/oscal",
				}),
			},
			wantImports: []oscalTypes.Import{
				{
					Href: "https://example.com/oscal",
					IncludeControls: &[]oscalTypes.SelectControlById{
						{
							WithIds: &[]string{
								"ex-1",
								"ex-1.2",
								"ex-2",
							},
						},
					},
				},
				{
					Href:       "testHref",
					IncludeAll: &oscalTypes.IncludeAll{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, profile, err := FromGuidance(&tt.guidance, "testHref", tt.options...)
			require.NoError(t, err)
			profileModel := oscalTypes.OscalModels{Profile: &profile}
			assert.NoError(t, oscalUtils.Validate(profileModel))

			assert.Equal(t, tt.wantImports, profile.Imports)
		})
	}
}

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
