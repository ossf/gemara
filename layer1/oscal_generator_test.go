package layer1

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

func TestToOSCALCatalog(t *testing.T) {
	tests := []struct {
		name       string
		guidance   GuidanceDocument
		wantGroups []oscalTypes.Group
		wantErr    bool
	}{
		{
			name:     "Good AIGF",
			guidance: goodAIGFExample(),
			wantGroups: []oscalTypes.Group{
				{
					Class: "category",
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
											Parts: &[]oscalTypes.Part{
												{
													Name: "guidance",
													ID:   "air-det-011_smt.1_gdn",
													Prose: "Define Intended Use and KPIs:\nObjectives: Clearly document how feedback data will be utilized, such as for prompt fine-tuning, RAG document updates,model/data drift detection, " +
														"or more advanced uses like Reinforcement Learning from Human Feedback (RLHF).\nKPI Alignment: Design feedback questions and metrics to align with the solution’s key performance indicators " +
														"(KPIs). For example, if accuracy is a KPI, feedback might involve users or SMEs annotating if an answer was correct.",
												},
											},
										},
										{
											Name:  "item",
											ID:    "air-det-011_smt.2",
											Title: "Types of Feedback and Collection Methods",
											Prose: "Implementing an effective human feedback loop involves clear collection processes.",
											Parts: &[]oscalTypes.Part{
												{
													Name: "guidance",
													ID:   "air-det-011_smt.2_gdn",
													Prose: "Quantitative Feedback:\nDescription: Involves collecting structured responses that can be easily aggregated and measured, such as numerical ratings (e.g., “Rate this response on " +
														"a scale of 1-5 for helpfulness”), categorical choices (e.g., “Was this answer: Correct/Incorrect/Partially Correct”), or binary responses (e.g., thumbs up/down)." +
														"\nUse Cases: Effective for tracking trends, measuring against KPIs, and quickly identifying areas of high or low performance.",
												},
											},
										},
									},
								},
								{
									Name: "assessment-objective",
									ID:   "air-det-011_obj",
									Prose: "A Human Feedback Loop is a critical detective and continuous improvement mechanism that involves systematically collecting, analyzing, and acting upon feedback provided by human users, " +
										"subject matter experts (SMEs), or reviewers regarding an AI system’s performance, outputs, or behavior.",
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
			guidance: GuidanceDocument{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			catalog, err := tt.guidance.ToOSCALCatalog()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				oscalDocument := oscalTypes.OscalModels{
					Catalog: &catalog,
				}
				err = oscalUtils.Validate(oscalDocument)
				assert.NoError(t, err)
				assert.Equal(t, tt.wantGroups, *catalog.Groups)
			}
		})
	}
}

func TestToOSCALProfile(t *testing.T) {
	guidanceWithImports := goodAIGFExample()
	// Add some shared guidelines
	mapping := MappingReference{
		Id:          "EXP",
		Description: "Example",
		Version:     "0.1.0",
		Url:         "https://example.com",
	}

	importedGuidelines := Mapping{
		ReferenceId: "EXP",
		Entries: []MappingEntry{
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
		guidance    GuidanceDocument
		options     []GenerateOption
		wantImports []oscalTypes.Import
	}{
		{
			name:     "Success/LocalOnly",
			guidance: goodAIGFExample(),
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
			profile, err := tt.guidance.ToOSCALProfile("testHref", tt.options...)
			require.NoError(t, err)
			oscalDocument := oscalTypes.OscalModels{
				Profile: &profile,
			}
			assert.NoError(t, oscalUtils.Validate(oscalDocument))

			assert.Equal(t, tt.wantImports, profile.Imports)
		})
	}
}
