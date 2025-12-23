package gemara

import (
	"fmt"
)

func ExamplePolicy_ToMarkdownChecklist() {
	policy := &Policy{
		Metadata: Metadata{
			Id:      "security-policy-v2.1",
			Version: "2.1.0",
			Author: Actor{
				Name:    "Security Team",
				Version: "1.0.0",
			},
		},
		Title: "Information Security Policy",
		Adherence: Adherence{
			AssessmentPlans: []AssessmentPlan{
				{
					Id:            "plan-access-control-monthly",
					RequirementId: "AC-01.01",
					Frequency:     "monthly",
					EvaluationMethods: []AcceptedMethod{
						{
							Type:        "automated",
							Description: "Run automated access control audit script",
						},
						{
							Type:        "manual",
							Description: "Review access control logs for anomalies",
						},
					},
					EvidenceRequirements: "Access logs, audit report, and review notes",
				},
				{
					Id:            "plan-encryption-quarterly",
					RequirementId: "SC-02.01",
					Frequency:     "quarterly",
					EvaluationMethods: []AcceptedMethod{
						{
							Type:        "automated",
							Description: "Verify encryption is enabled on all data stores",
						},
					},
					EvidenceRequirements: "Configuration scan results and compliance report",
				},
				{
					Id:            "plan-backup-automated",
					RequirementId: "CP-01.01",
					Frequency:     "weekly",
					EvaluationMethods: []AcceptedMethod{
						{
							Type:        "automated",
							Description: "Verify backup jobs completed successfully",
						},
					},
					EvidenceRequirements: "Backup job logs and success reports",
				},
				{
					Id:            "plan-backup-manual",
					RequirementId: "CP-01.01",
					Frequency:     "monthly",
					EvaluationMethods: []AcceptedMethod{
						{
							Type:        "manual",
							Description: "Test restore procedure from backup",
						},
					},
					EvidenceRequirements: "Restore test results and documentation",
				},
			},
		},
	}

	checklist, err := policy.ToMarkdownChecklist()
	if err != nil {
		fmt.Printf("Error generating checklist: %v\n", err)
		return
	}

	// Print the generated checklist
	fmt.Println(checklist)

	// Output:
	// # Policy Checklist: Information Security Policy (security-policy-v2.1)
	//
	// **Author:** Security Team (v1.0.0)
	//
	// ## Assessment Requirement: AC-01.01
	//
	// - [ ] Run automated access control audit script (monthly) [Plan: plan-access-control-monthly]
	//     > **Evidence Required:** Access logs, audit report, and review notes
	// - [ ] Review access control logs for anomalies (monthly) [Plan: plan-access-control-monthly]
	//     > **Evidence Required:** Access logs, audit report, and review notes
	//
	// ---
	//
	// ## Assessment Requirement: SC-02.01
	//
	// - [ ] Verify encryption is enabled on all data stores (quarterly) [Plan: plan-encryption-quarterly]
	//     > **Evidence Required:** Configuration scan results and compliance report
	//
	// ---
	//
	// ## Assessment Requirement: CP-01.01
	//
	// - [ ] Verify backup jobs completed successfully (weekly) [Plan: plan-backup-automated]
	//     > **Evidence Required:** Backup job logs and success reports
	// - [ ] Test restore procedure from backup (monthly) [Plan: plan-backup-manual]
	//     > **Evidence Required:** Restore test results and documentation
}
