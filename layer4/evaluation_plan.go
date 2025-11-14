package layer4

import (
	"bytes"
	"fmt"
	"text/template"
)

// ChecklistItem represents a single checklist item.
type ChecklistItem struct {
	RequirementId         string // Requirement ID (e.g., "AC-1.1")
	ProcedureName         string // Procedure name
	Description           string // Description (if different from name)
	Documentation         string // Documentation URL
	IsAdditionalProcedure bool   // Additional procedure
}

// ControlSection organizes checklist items by control.
type ControlSection struct {
	ControlName      string          // Control identifier (e.g., "AC-1")
	ControlReference string          // Formatted reference (e.g., "NIST-800-53 / AC-1")
	Items            []ChecklistItem // Checklist items for this control
}

// Checklist represents the structured checklist data.
type Checklist struct {
	PlanId        string           // Evaluation plan ID
	Author        string           // Author name
	AuthorVersion string           // Author version
	Sections      []ControlSection // Control sections
}

// ToChecklist converts an EvaluationPlan into a structured Checklist.
func (e EvaluationPlan) ToChecklist() Checklist {
	checklist := Checklist{}

	if e.Metadata.Id != "" {
		checklist.PlanId = e.Metadata.Id
	}
	if e.Metadata.Author.Name != "" {
		checklist.Author = e.Metadata.Author.Name
		checklist.AuthorVersion = e.Metadata.Author.Version
	}

	for _, plan := range e.Plans {
		if plan.Control.EntryId == "" {
			continue
		}

		// Get control name with fallback: EntryId -> ReferenceId -> default
		controlName := "Unnamed Control"
		if plan.Control.EntryId != "" {
			controlName = plan.Control.EntryId
		} else if plan.Control.ReferenceId != "" {
			controlName = plan.Control.ReferenceId
		}

		// Format control reference as "Framework / Control-ID" (e.g. OSPS-B / AC-01)
		controlReference := ""
		if plan.Control.ReferenceId != "" || plan.Control.EntryId != "" {
			controlReference = fmt.Sprintf("%s / %s", plan.Control.ReferenceId, plan.Control.EntryId)
		}

		section := ControlSection{
			ControlName:      controlName,
			ControlReference: controlReference,
			Items:            buildChecklistItems(&plan),
		}

		checklist.Sections = append(checklist.Sections, section)
	}

	return checklist
}

// ToMarkdownChecklist converts an evaluation plan into a markdown checklist.
// Generates a pre-execution checklist showing what needs to be checked.
func (e EvaluationPlan) ToMarkdownChecklist() string {
	checklist := e.ToChecklist()
	markdown, err := checklist.ToMarkdownWithTemplate(MarkdownTemplate)
	if err != nil {
		return ""
	}
	return markdown
}

// ToMarkdownWithTemplate formats the checklist using a custom template.
func (c Checklist) ToMarkdownWithTemplate(templateStr string) (string, error) {
	tmpl, err := template.New("checklist").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, c); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// buildChecklistItems converts an AssessmentPlan into checklist items.
func buildChecklistItems(plan *AssessmentPlan) []ChecklistItem {
	if plan == nil {
		return nil
	}

	if len(plan.Assessments) == 0 {
		return []ChecklistItem{
			{
				RequirementId: "",
				ProcedureName: "No assessments defined",
			},
		}
	}

	var items []ChecklistItem
	assessmentNum := 1

	for _, assessment := range plan.Assessments {
		requirementId := assessment.Requirement.EntryId
		if requirementId == "" {
			requirementId = fmt.Sprintf("Assessment %d", assessmentNum)
		}
		assessmentNum++

		if len(assessment.Procedures) == 0 {
			items = append(items, ChecklistItem{
				RequirementId: requirementId,
				ProcedureName: "No procedures defined",
			})
		} else {
			for i, procedure := range assessment.Procedures {
				// Get procedure name with fallback: Name -> Description -> Id
				procedureName := procedure.Id
				if procedure.Name != "" {
					procedureName = procedure.Name
				} else if procedure.Description != "" {
					procedureName = procedure.Description
				}

				item := ChecklistItem{
					RequirementId:         requirementId,
					ProcedureName:         procedureName,
					Description:           procedure.Description,
					Documentation:         procedure.Documentation,
					IsAdditionalProcedure: i > 0,
				}

				if i > 0 {
					item.RequirementId = ""
				}

				items = append(items, item)
			}
		}
	}

	return items
}
