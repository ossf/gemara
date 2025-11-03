package layer4

import (
	"fmt"
	"strings"
)

// ToMarkdownChecklist converts an evaluation plan into a markdown checklist for human-led
// assessments. Used BEFORE execution - shows what needs to be checked, not the results.
// Note: reference-id values must map to metadata.mapping-references[id], and requirement
// reference-ids must match their parent control's reference-id (enforced by schema).
func (e EvaluationPlan) ToMarkdownChecklist() string {
	var b strings.Builder

	if e.Metadata.Id != "" {
		fmt.Fprintf(&b, "# Evaluation Plan: %s\n\n", e.Metadata.Id)
	}
	if e.Metadata.Author.Name != "" {
		fmt.Fprintf(&b, "**Author:** %s", e.Metadata.Author.Name)
		if e.Metadata.Author.Version != "" {
			fmt.Fprintf(&b, " (v%s)", e.Metadata.Author.Version)
		}
		b.WriteString("\n\n")
	}

	isFirstPlan := true
	for _, plan := range e.Plans {
		if plan.Control.EntryId == "" {
			continue
		}

		// Add separator between controls (except before the first one)
		if !isFirstPlan {
			b.WriteString("\n---\n\n")
		}
		isFirstPlan = false

		controlChecklist := formatAssessmentPlanAsMarkdown(&plan)
		b.WriteString(controlChecklist)
	}

	return b.String()
}

// formatAssessmentPlanAsMarkdown creates a markdown representation of a single assessment plan.
// This is used internally by EvaluationPlan.ToMarkdownChecklist().
func formatAssessmentPlanAsMarkdown(plan *AssessmentPlan) string {
	if plan == nil {
		return ""
	}

	var b strings.Builder

	controlName := getControlName(plan)
	fmt.Fprintf(&b, "## %s\n\n", controlName)

	// Control ID and mapping info
	if plan.Control.ReferenceId != "" || plan.Control.EntryId != "" {
		fmt.Fprintf(&b, "**Control:** %s / %s\n\n", plan.Control.ReferenceId, plan.Control.EntryId)
	}

	if len(plan.Assessments) == 0 {
		b.WriteString("- [ ] No assessments defined\n")
	} else {
		assessmentNum := 1
		for _, assessment := range plan.Assessments {
			requirementId := assessment.Requirement.EntryId
			if requirementId == "" {
				requirementId = fmt.Sprintf("Assessment %d", assessmentNum)
			}
			assessmentNum++

			if len(assessment.Procedures) == 0 {
				fmt.Fprintf(&b, "- [ ] **%s**: No procedures defined\n", requirementId)
			} else {
				for i, procedure := range assessment.Procedures {
					procedureName := getProcedureName(&procedure)
					if i == 0 {
						// First procedure uses the requirement ID
						fmt.Fprintf(&b, "- [ ] **%s**: %s", requirementId, procedureName)
						if procedure.Description != "" && procedure.Description != procedureName {
							fmt.Fprintf(&b, " - %s", procedure.Description)
						}
						b.WriteString("\n")
					} else {
						// Additional procedures for the same requirement
						fmt.Fprintf(&b, "  - [ ] %s", procedureName)
						if procedure.Description != "" && procedure.Description != procedureName {
							fmt.Fprintf(&b, " - %s", procedure.Description)
						}
						b.WriteString("\n")
					}

					if procedure.Documentation != "" {
						fmt.Fprintf(&b, "    > [Documentation](%s)\n", procedure.Documentation)
					}
				}
			}
		}
	}

	return b.String()
}

// getControlName returns a human-readable name for the control from the assessment plan.
func getControlName(plan *AssessmentPlan) string {
	if plan == nil {
		return "Unknown Control"
	}
	// Use EntryId as the control name if no other name is available
	if plan.Control.EntryId != "" {
		return plan.Control.EntryId
	}
	if plan.Control.ReferenceId != "" {
		return plan.Control.ReferenceId
	}
	return "Unnamed Control"
}

// getProcedureName returns a human-readable name for a procedure, with fallback logic.
// Priority: Name -> Description -> Id
func getProcedureName(procedure *AssessmentProcedure) string {
	if procedure.Name != "" {
		return procedure.Name
	}
	if procedure.Description != "" {
		return procedure.Description
	}
	return procedure.Id
}
