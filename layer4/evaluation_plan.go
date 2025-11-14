package layer4

import (
	"bytes"
	"fmt"
	"text/template"
)

// ChecklistItem represents a single checklist item.
type ChecklistItem struct {
	// RequirementId is the requirement ID (e.g., "OSPS-AC-01.01")
	RequirementId string
	// ProcedureName is the human-readable name of the procedure to execute.
	ProcedureName string
	// Description provides additional context or a summary about the procedure.
	Description string
	// Documentation is the documentation URL
	Documentation string
	// IsAdditionalProcedure indicates if this is an additional procedure
	IsAdditionalProcedure bool
}

// ControlSection organizes checklist items by control.
type ControlSection struct {
	// ControlName is the control identifier (e.g., "OSPS-AC-01")
	ControlName string
	// ControlReference is the formatted reference (e.g., "OSPS-B / OSPS-AC-01")
	ControlReference string
	// Items are the checklist items for this control
	Items []ChecklistItem
}

// Checklist represents the structured checklist data.
type Checklist struct {
	// PlanId identifies the evaluation plan.
	PlanId string
	// Author is the name of the plan author.
	Author string
	// AuthorVersion is the version of the authoring tool or system.
	AuthorVersion string
	// Sections are the control sections
	Sections []ControlSection
}

// ToChecklist converts an EvaluationPlan into a structured Checklist.
func (e EvaluationPlan) ToChecklist() (Checklist, error) {
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

		// Format control reference as "Framework / Control-ID" (e.g. OSPS-B / OSPS-AC-01)
		controlReference := ""
		if plan.Control.ReferenceId != "" || plan.Control.EntryId != "" {
			controlReference = fmt.Sprintf("%s / %s", plan.Control.ReferenceId, plan.Control.EntryId)
		}

		items, err := buildChecklistItems(&plan)
		if err != nil {
			return Checklist{}, fmt.Errorf("failed to build checklist items for control %q: %w", controlName, err)
		}

		section := ControlSection{
			ControlName:      controlName,
			ControlReference: controlReference,
			Items:            items,
		}

		checklist.Sections = append(checklist.Sections, section)
	}

	return checklist, nil
}

// ToMarkdownChecklist converts an evaluation plan into a markdown checklist.
// Generates a pre-execution checklist showing what needs to be checked.
func (e EvaluationPlan) ToMarkdownChecklist() (string, error) {
	checklist, err := e.ToChecklist()
	if err != nil {
		return "", fmt.Errorf("failed to build checklist: %w", err)
	}

	tmpl, err := template.New("checklist").Parse(markdownTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, checklist); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// buildChecklistItems converts an AssessmentPlan into checklist items.
func buildChecklistItems(plan *AssessmentPlan) ([]ChecklistItem, error) {
	if plan == nil {
		return nil, fmt.Errorf("assessment plan is nil")
	}

	if len(plan.Assessments) == 0 {
		return nil, fmt.Errorf("assessment plan has no assessments")
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
			return nil, fmt.Errorf("assessment %q has no procedures", requirementId)
		}

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

			items = append(items, item)
		}
	}

	return items, nil
}
