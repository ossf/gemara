package checklist

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ossf/gemara"
	"github.com/stretchr/testify/require"
)

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot(t *testing.T) string {
	dir, err := os.Getwd()
	require.NoError(t, err)
	
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			t.Fatalf("Could not find project root (go.mod)")
		}
		dir = parent
	}
}

func Test_ToMarkdownChecklist(t *testing.T) {
	tests := []struct {
		name     string
		checklist Checklist
		contains []string
		notContains []string
	}{
		{
			name: "comprehensive checklist with multiple controls and requirements",
			checklist: Checklist{
				PolicyId: "policy-2024-01",
				Author:   "gemara",
				AuthorVersion: "1.0.0",
				Sections: []ControlSection{
					{
						ControlRef: "OSPS-B",
						Requirements: []RequirementSection{
							{
								RequirementId: "OSPS-AC-01.01",
								RequirementRecommendation: "Enforce multi-factor authentication for the project's version control system.",
								Items: []Item{
									{
										RequirementId:         "OSPS-AC-01.01",
										EvaluatorName:         "Security Assessment Team",
										Description:           "Check that MFA is configured for the repository",
										Documentation:         "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml",
										IsAdditionalEvaluator: false,
									},
									{
										RequirementId:         "OSPS-AC-01.01",
										EvaluatorName:         "OSPS Baseline Scanner",
										Description:           "Verify the policy contains required elements",
										Documentation:         "",
										IsAdditionalEvaluator: true,
									},
								},
							},
							{
								RequirementId: "OSPS-AC-01.02",
								RequirementRecommendation: "",
								Items: []Item{
									{
										RequirementId:         "OSPS-AC-01.02",
										EvaluatorName:         "Security Assessment Team",
										Description:           "Verify the policy has been approved by management",
										IsAdditionalEvaluator: false,
									},
								},
							},
						},
					},
					{
						ControlRef: "OSPS-B",
						Requirements: []RequirementSection{
							{
								RequirementId: "OSPS-AC-03.01",
								RequirementRecommendation: "Configure branch protection rules.",
								Items: []Item{
									{
										RequirementId:         "OSPS-AC-03.01",
										EvaluatorName:         "Security Assessment Team",
										Description:           "Check that the branch protection rules are configured for the primary branch",
										IsAdditionalEvaluator: false,
									},
								},
							},
						},
					},
				},
			},
			contains: []string{
				"# Evaluation Plan: policy-2024-01",
				"**Author:** gemara",
				"(v1.0.0)",
				"## OSPS-B",
				"### OSPS-AC-01.01",
				"Enforce multi-factor authentication",
				"- [ ] **Security Assessment Team** - Check that MFA is configured for the repository",
				"    > [Documentation](https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml)",
				"  - [ ]  - Verify the policy contains required elements",
				"### OSPS-AC-01.02",
				"- [ ] **Security Assessment Team** - Verify the policy has been approved by management",
				"---",
				"## OSPS-B",
				"### OSPS-AC-03.01",
				"Configure branch protection rules.",
			},
		},
		{
			name: "empty checklist",
			checklist: Checklist{
				PolicyId: "empty-policy",
				Author:   "test",
				Sections: []ControlSection{},
			},
			contains: []string{
				"# Evaluation Plan: empty-policy",
				"**Author:** test",
			},
			notContains: []string{
				"## Summary",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown, err := ToMarkdownChecklist(tt.checklist)
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
	policy := gemara.Policy{
		Metadata: gemara.Metadata{
			Id: "test-policy",
			Author: gemara.Actor{
				Name:    "test-author",
				Version: "1.0.0",
			},
		},
		ImplementationPlan: gemara.ImplementationPlan{
			Evaluators: []gemara.Actor{
				{
					Id:          "security-team",
					Name:        "Security Assessment Team",
					Description: "Check that MFA is configured for the repository",
					Uri:         "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml",
				},
			},
		},
		ControlReferences: []gemara.PolicyMapping{
			{
				ReferenceId: "OSPS-B",
				AssessmentRequirementModifications: []gemara.AssessmentRequirementModifier{
					{
						TargetId: "OSPS-AC-01.01",
						Extensions: &gemara.AssessmentRequirementExtensions{
							RequiredEvaluators: []string{"security-team"},
						},
					},
				},
			},
		},
	}

	checklist, err := ToChecklist(policy)
	require.NoError(t, err)

	require.Equal(t, "test-policy", checklist.PolicyId)
	require.Equal(t, "test-author", checklist.Author)
	require.Equal(t, "1.0.0", checklist.AuthorVersion)
	require.Len(t, checklist.Sections, 1)

	section := checklist.Sections[0]
	require.Equal(t, "OSPS-B", section.ControlRef)
	require.Len(t, section.Requirements, 1)

	requirement := section.Requirements[0]
	require.Equal(t, "OSPS-AC-01.01", requirement.RequirementId)
	require.Len(t, requirement.Items, 1)

	item := requirement.Items[0]
	require.Equal(t, "OSPS-AC-01.01", item.RequirementId)
	require.Equal(t, "Security Assessment Team", item.EvaluatorName)
	require.Equal(t, "Check that MFA is configured for the repository", item.Description)
	require.Equal(t, "https://github.com/ossf/security-baseline/blob/main/baseline/OSPS-AC.yaml", item.Documentation)
	require.False(t, item.IsAdditionalEvaluator)
}

func Test_ToChecklist_ErrorCases(t *testing.T) {
	t.Run("missing evaluator in map", func(t *testing.T) {
		policy := gemara.Policy{
			Metadata: gemara.Metadata{
				Id: "test-policy",
			},
			ImplementationPlan: gemara.ImplementationPlan{
				Evaluators: []gemara.Actor{
					{
						Id:   "other-evaluator",
						Name: "Other Evaluator",
					},
				},
			},
			ControlReferences: []gemara.PolicyMapping{
				{
					ReferenceId: "OSPS-B",
					AssessmentRequirementModifications: []gemara.AssessmentRequirementModifier{
						{
							TargetId: "OSPS-AC-01.01",
							Extensions: &gemara.AssessmentRequirementExtensions{
								RequiredEvaluators: []string{"non-existent-evaluator"},
							},
						},
					},
				},
			},
		}

		// This should succeed but the evaluator will have empty name/description
		checklist, err := ToChecklist(policy)
		require.NoError(t, err)
		require.Len(t, checklist.Sections, 1)
		require.Len(t, checklist.Sections[0].Requirements, 1)
		// The evaluator won't be found, so it will have empty name
		require.Len(t, checklist.Sections[0].Requirements[0].Items, 1)
		require.Equal(t, "", checklist.Sections[0].Requirements[0].Items[0].EvaluatorName)
	})

	t.Run("no required evaluators", func(t *testing.T) {
		policy := gemara.Policy{
			Metadata: gemara.Metadata{
				Id: "test-policy",
			},
			ImplementationPlan: gemara.ImplementationPlan{
				Evaluators: []gemara.Actor{
					{
						Id:   "security-team",
						Name: "Security Assessment Team",
					},
				},
			},
			ControlReferences: []gemara.PolicyMapping{
				{
					ReferenceId: "OSPS-B",
					AssessmentRequirementModifications: []gemara.AssessmentRequirementModifier{
						{
							TargetId:   "OSPS-AC-01.01",
							Extensions: &gemara.AssessmentRequirementExtensions{
								RequiredEvaluators: []string{},
							},
						},
					},
				},
			},
		}

		_, err := ToChecklist(policy)
		require.Error(t, err)
		require.Contains(t, err.Error(), "has no evaluators")
	})
}

func Test_ToChecklist_WithCatalog(t *testing.T) {
	policy := gemara.Policy{
		Metadata: gemara.Metadata{
			Id: "test-policy",
			Author: gemara.Actor{
				Name:    "test-author",
				Version: "1.0.0",
			},
		},
		ImplementationPlan: gemara.ImplementationPlan{
			Evaluators: []gemara.Actor{
				{
					Id:          "security-team",
					Name:        "Security Assessment Team",
					Description: "Check MFA configuration",
				},
			},
		},
		ControlReferences: []gemara.PolicyMapping{
			{
				ReferenceId: "OSPS-B",
				AssessmentRequirementModifications: []gemara.AssessmentRequirementModifier{
					{
						TargetId: "OSPS-AC-01.01",
						Extensions: &gemara.AssessmentRequirementExtensions{
							RequiredEvaluators: []string{"security-team"},
						},
					},
				},
			},
		},
	}

	// Add mapping reference to policy metadata pointing to test catalog
	// Resolve absolute path relative to project root in test to ensure it works from any directory
	projectRoot := findProjectRoot(t)
	testDataPath := filepath.Join(projectRoot, "test-data", "good-osps.yml")
	absPath, err := filepath.Abs(testDataPath)
	require.NoError(t, err, "Failed to resolve test data path")
	
	policy.Metadata.MappingReferences = []gemara.MappingReference{
		{
			Id:  "OSPS-B",
			Url: "file://" + absPath,
		},
	}

	checklist, err := ToChecklist(policy)
	require.NoError(t, err)

	require.Equal(t, "test-policy", checklist.PolicyId)
	require.Len(t, checklist.Sections, 1)

	section := checklist.Sections[0]
	require.Len(t, section.Requirements, 1)

	requirement := section.Requirements[0]
	require.Equal(t, "OSPS-AC-01.01", requirement.RequirementId)
	// Should have description from Layer 2 catalog since recommendation is not overridden
	require.NotEmpty(t, requirement.RequirementRecommendation)
	require.Contains(t, requirement.RequirementRecommendation, "multi-factor authentication")
	require.Len(t, requirement.Items, 1)
}

func Test_ToChecklist_WithOverride(t *testing.T) {
	policy := gemara.Policy{
		Metadata: gemara.Metadata{
			Id: "test-policy",
		},
		ImplementationPlan: gemara.ImplementationPlan{
			Evaluators: []gemara.Actor{
				{
					Id:   "security-team",
					Name: "Security Assessment Team",
				},
			},
		},
		ControlReferences: []gemara.PolicyMapping{
			{
				ReferenceId: "OSPS-B",
				AssessmentRequirementModifications: []gemara.AssessmentRequirementModifier{
					{
						TargetId: "OSPS-AC-01.01",
						Overrides: &gemara.AssessmentRequirement{
							Id:   "OSPS-AC-01.01",
							Text: "Overridden text",
							Recommendation: "Overridden recommendation",
						},
						Extensions: &gemara.AssessmentRequirementExtensions{
							RequiredEvaluators: []string{"security-team"},
						},
					},
				},
			},
		},
	}

	// Add mapping reference to policy metadata (not needed for override test, but good to have)
	// Resolve absolute path relative to project root in test to ensure it works from any directory
	projectRoot := findProjectRoot(t)
	testDataPath := filepath.Join(projectRoot, "test-data", "good-osps.yml")
	absPath, err := filepath.Abs(testDataPath)
	require.NoError(t, err, "Failed to resolve test data path")
	
	policy.Metadata.MappingReferences = []gemara.MappingReference{
		{
			Id:  "OSPS-B",
			Url: "file://" + absPath,
		},
	}

	// With override, should use override text regardless of catalog availability
	checklist, err := ToChecklist(policy)
	require.NoError(t, err)

	requirement := checklist.Sections[0].Requirements[0]
	// Should use override text, not Layer 2 text
	require.Equal(t, "Overridden text", requirement.RequirementRecommendation)
}
