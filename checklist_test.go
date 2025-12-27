package gemara

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPolicy_ToMarkdownChecklist(t *testing.T) {
	tests := []struct {
		name     string
		policy   *Policy
		wantErr  bool
		validate func(*testing.T, string)
	}{
		{
			name: "Single Plan",
			policy: &Policy{
				Metadata: Metadata{
					Id: "test-policy",
					Author: Actor{
						Name:    "Test Author",
						Version: "1.0.0",
					},
				},
				Title: "Test Security Policy",
				Adherence: Adherence{
					AssessmentPlans: []AssessmentPlan{
						{
							Id:            "plan-1",
							RequirementId: "REQ-001",
							Frequency:     "quarterly",
							EvaluationMethods: []AcceptedMethod{
								{
									Type:        "automated",
									Description: "Run automated security scan",
								},
							},
							EvidenceRequirements: "Scan report",
						},
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, got string) {
				assert.Contains(t, got, "Test Security Policy")
				assert.Contains(t, got, "[Plan: plan-1]")
				assert.Contains(t, got, "Run automated security scan")
			},
		},
		{
			name: "Multiple Plans",
			policy: &Policy{
				Metadata: Metadata{
					Id: "multi-plan-policy",
					Author: Actor{
						Name: "Multi Author",
					},
				},
				Title: "Multi Plan Policy",
				Adherence: Adherence{
					AssessmentPlans: []AssessmentPlan{
						{
							Id:            "plan-1",
							RequirementId: "REQ-001",
							Frequency:     "monthly",
							EvaluationMethods: []AcceptedMethod{
								{Type: "automated", Description: "Automated check"},
							},
						},
						{
							Id:            "plan-2",
							RequirementId: "REQ-001",
							Frequency:     "quarterly",
							EvaluationMethods: []AcceptedMethod{
								{Type: "manual", Description: "Manual review"},
							},
						},
						{
							Id:            "plan-3",
							RequirementId: "REQ-002",
							Frequency:     "annually",
							EvaluationMethods: []AcceptedMethod{
								{Type: "automated"},
							},
						},
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, got string) {
				assert.Equal(t, 1, strings.Count(got, "## Assessment Requirement: REQ-001"))
				assert.Equal(t, 1, strings.Count(got, "## Assessment Requirement: REQ-002"))
				assert.Contains(t, got, "[Plan: plan-1]")
				assert.Contains(t, got, "[Plan: plan-2]")
			},
		},
		{
			name: "Multiple Methods",
			policy: &Policy{
				Metadata: Metadata{Id: "multi-method-policy"},
				Title:    "Multi Method Policy",
				Adherence: Adherence{
					AssessmentPlans: []AssessmentPlan{
						{
							Id:            "plan-1",
							RequirementId: "REQ-001",
							Frequency:     "monthly",
							EvaluationMethods: []AcceptedMethod{
								{
									Type:        "automated",
									Description: "Primary automated check",
								},
								{
									Type:        "manual",
									Description: "Secondary manual review",
								},
							},
						},
					},
				},
			},
			wantErr: false,
			validate: func(t *testing.T, got string) {
				assert.Contains(t, got, "Primary automated check")
				assert.Contains(t, got, "Secondary manual review")
			},
		},
		{
			name: "No Methods Error",
			policy: &Policy{
				Metadata: Metadata{Id: "no-methods-policy"},
				Title:    "No Methods Policy",
				Adherence: Adherence{
					AssessmentPlans: []AssessmentPlan{
						{
							Id:                "plan-1",
							RequirementId:     "REQ-001",
							Frequency:         "quarterly",
							EvaluationMethods: []AcceptedMethod{},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.policy.ToMarkdownChecklist()
			if tt.wantErr {
				require.Error(t, err, "expected error but got none")
				return
			}
			require.NoError(t, err, "unexpected error generating checklist")
			if tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}
