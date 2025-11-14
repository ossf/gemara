package layer4

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToMarkdownChecklist(t *testing.T) {
	tests := []struct {
		name           string
		evaluationPlan EvaluationPlan
		contains       []string
		notContains    []string
	}{
		{
			name: "comprehensive evaluation plan with multiple controls and assessments",
			evaluationPlan: EvaluationPlan{
				Plans: []AssessmentPlan{
					{
						Control: Mapping{
							ReferenceId: "NIST-800-53",
							EntryId:     "AC-1",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "NIST-800-53",
									EntryId:     "AC-1.1",
								},
								Procedures: []AssessmentProcedure{
									{
										Id:            "proc-1",
										Name:          "Verify access control policy",
										Description:   "Check that an access control policy exists and is documented",
										Documentation: "https://example.com/docs/ac-1",
									},
									{
										Id:          "proc-2",
										Name:        "Review policy content",
										Description: "Verify the policy contains required elements",
									},
								},
							},
							{
								Requirement: Mapping{
									ReferenceId: "NIST-800-53",
									EntryId:     "AC-1.2",
								},
								Procedures: []AssessmentProcedure{
									{
										Id:          "proc-3",
										Name:        "Check policy approval",
										Description: "Verify the policy has been approved by management",
									},
								},
							},
						},
					},
					{
						Control: Mapping{
							ReferenceId: "CIS",
							EntryId:     "CTRL-2",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "CIS",
									EntryId:     "CTRL-2.1",
								},
								Procedures: []AssessmentProcedure{
									{
										Name:        "Verify configuration setting",
										Description: "Check that the configuration is properly set",
									},
								},
							},
						},
					},
				},
				Metadata: Metadata{
					Id:      "plan-2024-01",
					Version: "1.0.0",
					Author: Author{
						Name:    "gemara",
						Uri:     "https://github.com/ossf/gemara",
						Version: "1.0.0",
					},
					MappingReferences: []MappingReference{
						{
							Id:      "NIST-800-53",
							Title:   "NIST Special Publication 800-53",
							Version: "Rev. 5",
							Url:     "https://csrc.nist.gov/publications/detail/sp/800-53/rev-5/final",
						},
						{
							Id:      "CIS",
							Title:   "CIS Controls",
							Version: "v8",
							Url:     "https://www.cisecurity.org/controls",
						},
					},
				},
			},
			contains: []string{
				"# Evaluation Plan: plan-2024-01",
				"**Author:** gemara",
				"(v1.0.0)",
				"## AC-1",
				"**Control:** NIST-800-53 / AC-1",
				"- [ ] **AC-1.1**: Verify access control policy - Check that an access control policy exists and is documented",
				"    > [Documentation](https://example.com/docs/ac-1)",
				"  - [ ] Review policy content - Verify the policy contains required elements",
				"- [ ] **AC-1.2**: Check policy approval - Verify the policy has been approved by management",
				"---",
				"## CTRL-2",
				"**Control:** CIS / CTRL-2",
				"- [ ] **CTRL-2.1**: Verify configuration setting - Check that the configuration is properly set",
			},
		},
		{
			name: "edge cases: empty plan, missing names, and empty IDs",
			evaluationPlan: EvaluationPlan{
				Plans: []AssessmentPlan{
					{
						Control: Mapping{
							ReferenceId: "TEST-REF",
							EntryId:     "EDGE-CASE-CTRL",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "TEST-REF",
									EntryId:     "",
								}, // Empty ID should be numbered
								Procedures: []AssessmentProcedure{
									{
										// No name or description - should use numbered fallback
										Id: "proc-1",
									},
								},
							},
							{
								Requirement: Mapping{
									ReferenceId: "TEST-REF",
									EntryId:     "VALID-ID",
								},
								Procedures: []AssessmentProcedure{
									{
										Name: "Test procedure",
									},
								},
							},
						},
					},
				},
				Metadata: Metadata{
					Author: Author{Name: "test"},
					MappingReferences: []MappingReference{
						{
							Id:      "TEST-REF",
							Title:   "Test Reference",
							Version: "1.0",
						},
					},
				},
			},
			contains: []string{
				"**Author:** test",
				"## EDGE-CASE-CTRL",
				"**Control:** TEST-REF / EDGE-CASE-CTRL",
				"- [ ] **Assessment 1**: proc-1",
				"- [ ] **VALID-ID**: Test procedure",
			},
		},
		{
			name: "empty evaluation plan",
			evaluationPlan: EvaluationPlan{
				Plans: []AssessmentPlan{},
				Metadata: Metadata{
					Id: "empty-plan",
					Author: Author{
						Name: "test",
					},
				},
			},
			contains: []string{
				"# Evaluation Plan: empty-plan",
				"**Author:** test",
			},
			notContains: []string{
				"## Summary",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown := tt.evaluationPlan.ToMarkdownChecklist()
			require.NotEmpty(t, markdown)

			for _, expected := range tt.contains {
				require.Contains(t, markdown, expected,
					"Markdown should contain: %s", expected)
			}

			for _, notExpected := range tt.notContains {
				require.NotContains(t, markdown, notExpected,
					"Markdown should not contain: %s", notExpected)
			}
		})
	}
}

func Test_ToChecklist(t *testing.T) {
	plan := EvaluationPlan{
		Plans: []AssessmentPlan{
			{
				Control: Mapping{
					ReferenceId: "NIST-800-53",
					EntryId:     "AC-1",
				},
				Assessments: []Assessment{
					{
						Requirement: Mapping{
							ReferenceId: "NIST-800-53",
							EntryId:     "AC-1.1",
						},
						Procedures: []AssessmentProcedure{
							{
								Name:          "Test procedure",
								Description:   "Test description",
								Documentation: "https://example.com",
							},
						},
					},
				},
			},
		},
		Metadata: Metadata{
			Id: "test-plan",
			Author: Author{
				Name:    "test-author",
				Version: "1.0.0",
			},
		},
	}

	checklist := plan.ToChecklist()

	require.Equal(t, "test-plan", checklist.PlanId)
	require.Equal(t, "test-author", checklist.Author)
	require.Equal(t, "1.0.0", checklist.AuthorVersion)
	require.Len(t, checklist.Sections, 1)

	section := checklist.Sections[0]
	require.Equal(t, "AC-1", section.ControlName)
	require.Equal(t, "NIST-800-53 / AC-1", section.ControlReference)
	require.Len(t, section.Items, 1)

	item := section.Items[0]
	require.Equal(t, "AC-1.1", item.RequirementId)
	require.Equal(t, "Test procedure", item.ProcedureName)
	require.Equal(t, "Test description", item.Description)
	require.Equal(t, "https://example.com", item.Documentation)
	require.False(t, item.IsAdditionalProcedure)
}
