package oscal

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
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
		assertFunc func(*testing.T, oscalTypes.Catalog)
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
		{
			name:     "Success/MultiLevelNestedControls",
			guidance: guidanceWithMultiLevelNested(),
			wantErr:  false,
			assertFunc: func(t *testing.T, catalog oscalTypes.Catalog) {
				require.NotNil(t, catalog.Groups)
				groups := *catalog.Groups
				require.Len(t, groups, 1)

				acGroup := groups[0]
				assert.Equal(t, "AC", acGroup.ID)
				require.NotNil(t, acGroup.Controls)
				controls := *acGroup.Controls
				require.Len(t, controls, 1)

				ac1 := controls[0]
				assert.Equal(t, "ac-1", ac1.ID)
				require.NotNil(t, ac1.Controls)
				ac1Children := *ac1.Controls
				require.Len(t, ac1Children, 1)

				ac1Enh := ac1Children[0]
				assert.Equal(t, "ac-1-enh", ac1Enh.ID)
				require.NotNil(t, ac1Enh.Controls)
				ac1EnhChildren := *ac1Enh.Controls
				require.Len(t, ac1EnhChildren, 1)

				ac1Enh2 := ac1EnhChildren[0]
				assert.Equal(t, "ac-1-enh-2", ac1Enh2.ID)
				assert.Nil(t, ac1Enh2.Controls)
			},
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
				if tt.assertFunc != nil {
					tt.assertFunc(t, catalog)
				} else {
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
			}
		})
	}
}

func TestProfileFromGuidanceDocument(t *testing.T) {
	goodAIFG, err := goodAIGFExample()
	require.NoError(t, err)

	guidanceWithImports := guidanceWithImports(goodAIFG)
	guidanceWithExternalExtends := guidanceWithExternalExtends()
	guidanceWithMerging := guidanceWithMerging()
	guidanceWithLocalExtends := guidanceWithLocalExtends()

	tests := []struct {
		name               string
		guidance           gemara.GuidanceDocument
		options            []GenerateOption
		wantModify         bool
		wantAlterations    int
		wantImports        int
		wantImportControls bool
		assertFunc         func(*testing.T, oscalTypes.Profile)
	}{
		{
			name:               "Success/LocalOnly",
			guidance:           goodAIFG,
			wantModify:         false,
			wantImports:        1,
			wantImportControls: false,
		},
		{
			name:               "Success/WithImports",
			guidance:           guidanceWithImports,
			wantModify:         false,
			wantImports:        1,
			wantImportControls: false,
		},
		{
			name:     "Success/WithImportOverride",
			guidance: guidanceWithImports,
			options: []GenerateOption{
				WithOSCALImports(map[string]string{
					"EXP": "https://example.com/oscal",
				}),
			},
			wantModify:         false,
			wantImports:        1,
			wantImportControls: false,
		},
		{
			name:     "Success/WithExternalExtends",
			guidance: guidanceWithExternalExtends,
			options: []GenerateOption{
				WithOSCALImports(map[string]string{
					"NIST-800-53": "https://nist.gov/800-53",
				}),
			},
			wantModify:         true,
			wantAlterations:    1,
			wantImports:        2,
			wantImportControls: true,
		},
		{
			name:     "Success/WithMerging",
			guidance: guidanceWithMerging,
			options: []GenerateOption{
				WithOSCALImports(map[string]string{
					"NIST-800-53": "https://nist.gov/800-53",
				}),
			},
			wantModify:         true,
			wantAlterations:    1,
			wantImports:        2,
			wantImportControls: true,
			assertFunc:         assertMergedParts,
		},
		{
			name:               "Success/WithLocalExtends",
			guidance:           guidanceWithLocalExtends,
			wantModify:         false,
			wantImports:        1,
			wantImportControls: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, profile, err := FromGuidance(&tt.guidance, "testHref", tt.options...)
			require.NoError(t, err)
			oscalDocument := oscalTypes.OscalModels{
				Profile: &profile,
			}
			assert.NoError(t, oscalUtils.Validate(oscalDocument))

			if tt.wantModify {
				require.NotNil(t, profile.Modify)
				require.NotNil(t, profile.Modify.Alters)
				assert.Equal(t, tt.wantAlterations, len(*profile.Modify.Alters))

				alterations := *profile.Modify.Alters
				assert.Equal(t, "ac-1", alterations[0].ControlId)
				require.NotNil(t, alterations[0].Adds)
				assert.Greater(t, len(*alterations[0].Adds), 0)
			} else {
				assert.Nil(t, profile.Modify)
			}

			assert.Equal(t, tt.wantImports, len(profile.Imports))

			hasControlImport := false
			for _, imp := range profile.Imports {
				if imp.IncludeControls != nil {
					hasControlImport = true
					selectors := *imp.IncludeControls
					assert.Greater(t, len(selectors), 0)
					if len(selectors) > 0 {
						assert.NotNil(t, selectors[0].WithIds)
						controlIds := *selectors[0].WithIds
						assert.Contains(t, controlIds, "ac-1")
					}
				}
			}
			assert.Equal(t, tt.wantImportControls, hasControlImport)

			if tt.assertFunc != nil {
				tt.assertFunc(t, profile)
			}
		})
	}
}

func assertMergedParts(t *testing.T, profile oscalTypes.Profile) {
	require.NotNil(t, profile.Modify)
	require.NotNil(t, profile.Modify.Alters)
	alterations := *profile.Modify.Alters
	require.NotNil(t, alterations[0].Adds)
	require.Greater(t, len(*alterations[0].Adds), 0)

	firstAddition := (*alterations[0].Adds)[0]
	require.NotNil(t, firstAddition.Parts)
	parts := *firstAddition.Parts

	assert.GreaterOrEqual(t, len(parts), 2, "Merged parts should have parts from both guidelines")
	hasFirstStatement := false
	hasSecondStatement := false
	for _, part := range parts {
		if part.Name == "statement" && part.Parts != nil {
			for _, subPart := range *part.Parts {
				if subPart.Prose == "First statement" {
					hasFirstStatement = true
				}
				if subPart.Prose == "Second statement" {
					hasSecondStatement = true
				}
			}
		}
	}
	assert.True(t, hasFirstStatement, "First statement should be present in merged parts")
	assert.True(t, hasSecondStatement, "Second statement should be present in merged parts")
}
