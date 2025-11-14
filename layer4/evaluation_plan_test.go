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
							ReferenceId: "OSPS-B",
							EntryId:     "OSPS-AC-01",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "OSPS-B",
									EntryId:     "OSPS-AC-01.01",
								},
								Procedures: []AssessmentProcedure{
									{
										Id:            "test_multi_factor_authentication",
										Name:          "Verify MFA configured for repository",
										Description:   "Check that MFA is configured for the repository",
										Documentation: "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml",
									},
									{
										Id:          "test_review_policy_content",
										Name:        "Review policy content",
										Description: "Verify the policy contains required elements",
									},
								},
							},
							{
								Requirement: Mapping{
									ReferenceId: "OSPS-B",
									EntryId:     "OSPS-AC-01.02",
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
							ReferenceId: "OSPS-B",
							EntryId:     "OSPS-AC-03",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "OSPS-B",
									EntryId:     "OSPS-AC-03.01",
								},
								Procedures: []AssessmentProcedure{
									{
										Id:          "github_branch_protection",
										Name:        "Primary Branch Protection Requirements",
										Description: "Check that the branch protection rules are configured for the primary branch",
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
							Id:      "OSPS-B",
							Title:   "Open Source Project Security Baseline",
							Version: "1.0",
							Url:     "https://github.com/ossf/security-baseline/tree/main/baseline",
						},
					},
				},
			},
			contains: []string{
				"# Evaluation Plan: plan-2024-01",
				"**Author:** gemara",
				"(v1.0.0)",
				"## OSPS-AC-01",
				"**Control:** OSPS-B / OSPS-AC-01",
				"- [ ] **OSPS-AC-01.01**: Verify MFA configured for repository - Check that MFA is configured for the repository",
				"    > [Documentation](https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml)",
				"  - [ ] Review policy content - Verify the policy contains required elements",
				"- [ ] **OSPS-AC-01.02**: Check policy approval - Verify the policy has been approved by management",
				"---",
				"## OSPS-AC-03",
				"**Control:** OSPS-B / OSPS-AC-03",
				"- [ ] **OSPS-AC-03.01**: Primary Branch Protection Requirements - Check that the branch protection rules are configured for the primary branch",
			},
		},
		{
			name: "edge cases: empty plan, missing names, and empty IDs",
			evaluationPlan: EvaluationPlan{
				Plans: []AssessmentPlan{
					{
						Control: Mapping{
							ReferenceId: "OSPS-B",
							EntryId:     "OSPS-AC-02",
						},
						Assessments: []Assessment{
							{
								Requirement: Mapping{
									ReferenceId: "OSPS-B",
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
									ReferenceId: "OSPS-B",
									EntryId:     "OSPS-AC-02.02",
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
							Id:      "OSPS-B",
							Title:   "OSPS Baseline",
							Version: "1.0",
						},
					},
				},
			},
			contains: []string{
				"**Author:** test",
				"## OSPS-AC-02",
				"**Control:** OSPS-B / OSPS-AC-02",
				"- [ ] **Assessment 1**: proc-1",
				"- [ ] **OSPS-AC-02.02**: Test procedure",
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
			markdown, err := tt.evaluationPlan.ToMarkdownChecklist()
			require.NoError(t, err)
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
					ReferenceId: "OSPS-B",
					EntryId:     "OSPS-AC-01",
				},
				Assessments: []Assessment{
					{
						Requirement: Mapping{
							ReferenceId: "OSPS-B",
							EntryId:     "OSPS-AC-01.01",
						},
						Procedures: []AssessmentProcedure{
							{
								Name:          "Verify MFA configured for repository",
								Description:   "Check that MFA is configured for the repository",
								Documentation: "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml",
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

	checklist, err := plan.ToChecklist()
	require.NoError(t, err)

	require.Equal(t, "test-plan", checklist.PlanId)
	require.Equal(t, "test-author", checklist.Author)
	require.Equal(t, "1.0.0", checklist.AuthorVersion)
	require.Len(t, checklist.Sections, 1)

	section := checklist.Sections[0]
	require.Equal(t, "OSPS-AC-01", section.ControlName)
	require.Equal(t, "OSPS-B / OSPS-AC-01", section.ControlReference)
	require.Len(t, section.Items, 1)

	item := section.Items[0]
	require.Equal(t, "OSPS-AC-01.01", item.RequirementId)
	require.Equal(t, "Verify MFA configured for repository", item.ProcedureName)
	require.Equal(t, "Check that MFA is configured for the repository", item.Description)
	require.Equal(t, "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml", item.Documentation)
	require.False(t, item.IsAdditionalProcedure)
}

func Test_ToChecklist_ErrorCases(t *testing.T) {
	t.Run("no assessments", func(t *testing.T) {
		plan := EvaluationPlan{
			Plans: []AssessmentPlan{
				{
					Control: Mapping{
						ReferenceId: "OSPS-B",
						EntryId:     "OSPS-AC-01",
					},
					Assessments: []Assessment{},
				},
			},
		}

		_, err := plan.ToChecklist()
		require.Error(t, err)
		require.Contains(t, err.Error(), "has no assessments")
	})

	t.Run("no procedures", func(t *testing.T) {
		plan := EvaluationPlan{
			Plans: []AssessmentPlan{
				{
					Control: Mapping{
						ReferenceId: "OSPS-B",
						EntryId:     "OSPS-AC-01",
					},
					Assessments: []Assessment{
						{
							Requirement: Mapping{
								ReferenceId: "OSPS-B",
								EntryId:     "OSPS-AC-01.01",
							},
							Procedures: []AssessmentProcedure{},
						},
					},
				},
			},
		}

		_, err := plan.ToChecklist()
		require.Error(t, err)
		require.Contains(t, err.Error(), "has no procedures")
	})
}
