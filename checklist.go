package gemara

import (
	"bytes"
	"fmt"
	"text/template"
)

// ChecklistItem represents a single checklist item.
type ChecklistItem struct {
	// PlanId is the assessment plan identifier this item belongs to.
	PlanId string
	// MethodDescription provides additional context or a summary about the method.
	MethodDescription string
	// MethodType defines the category the method falls into.
	MethodType string
	// Frequency indicates how often this assessment should be performed.
	Frequency string
	// EvidenceRequirements describes what evidence is required for this assessment.
	EvidenceRequirements string
}

// RequirementSection organizes checklist items by assessment requirement.
type RequirementSection struct {
	// RequirementId is the assessment requirement identifier (e.g., "OSPS-AC-01.01")
	RequirementId string
	// Items are the checklist items for this requirement
	Items []ChecklistItem
}

// Checklist represents the structured checklist data.
type Checklist struct {
	// PolicyId identifies the policy.
	PolicyId string
	// PolicyTitle is the title of the policy.
	PolicyTitle string
	// Author is the name of the policy author.
	Author string
	// AuthorVersion is the version of the authoring tool or system.
	AuthorVersion string
	// Sections are the requirement sections
	Sections []RequirementSection
}

// ToMarkdownChecklist converts a policy into a markdown checklist.
func (p *Policy) ToMarkdownChecklist() (string, error) {
	checklist, err := p.toChecklist()
	if err != nil {
		return "", fmt.Errorf("failed to build checklist: %w", err)
	}

	tmpl, err := template.New("checklist").Parse(markdownChecklistTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, checklist); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// toChecklist converts a Policy into a structured Checklist.
func (p *Policy) toChecklist() (Checklist, error) {
	checklist := Checklist{}

	if p.Metadata.Id != "" {
		checklist.PolicyId = p.Metadata.Id
	}
	if p.Title != "" {
		checklist.PolicyTitle = p.Title
	}
	if p.Metadata.Author.Name != "" {
		checklist.Author = p.Metadata.Author.Name
		if p.Metadata.Author.Version != "" {
			checklist.AuthorVersion = p.Metadata.Author.Version
		}
	}

	if len(p.Adherence.AssessmentPlans) == 0 {
		return checklist, nil
	}

	// Group assessment plans by requirement-id
	requirementMap := make(map[string][]AssessmentPlan)
	var requirementOrder []string
	for _, plan := range p.Adherence.AssessmentPlans {
		if plan.RequirementId == "" {
			continue
		}
		if _, exists := requirementMap[plan.RequirementId]; !exists {
			requirementOrder = append(requirementOrder, plan.RequirementId)
		}
		requirementMap[plan.RequirementId] = append(requirementMap[plan.RequirementId], plan)
	}

	// Build sections from grouped plans in insertion order
	for _, requirementId := range requirementOrder {
		plans := requirementMap[requirementId]
		var allItems []ChecklistItem

		for _, plan := range plans {
			items, err := buildChecklistItems(&plan)
			if err != nil {
				return Checklist{}, fmt.Errorf("failed to build checklist items for requirement %q: %w", requirementId, err)
			}
			allItems = append(allItems, items...)
		}

		if len(allItems) == 0 {
			continue
		}

		section := RequirementSection{
			RequirementId: requirementId,
			Items:         allItems,
		}

		checklist.Sections = append(checklist.Sections, section)
	}

	return checklist, nil
}

// buildChecklistItems converts an AssessmentPlan into checklist items.
// Each evaluation method becomes a checklist item.
func buildChecklistItems(plan *AssessmentPlan) ([]ChecklistItem, error) {
	if plan == nil {
		return nil, fmt.Errorf("assessment plan is nil")
	}

	if len(plan.EvaluationMethods) == 0 {
		return nil, fmt.Errorf("assessment plan %q has no evaluation methods", plan.Id)
	}

	var items []ChecklistItem

	for _, method := range plan.EvaluationMethods {
		item := ChecklistItem{
			PlanId:               plan.Id,
			MethodDescription:    method.Description,
			MethodType:           method.Type,
			Frequency:            plan.Frequency,
			EvidenceRequirements: plan.EvidenceRequirements,
		}

		items = append(items, item)
	}

	return items, nil
}

// markdownChecklistTemplate is the default template for generating markdown checklist output.
// This template is used internally by ToMarkdownChecklist().
const markdownChecklistTemplate = `{{if .PolicyId}}# Policy Checklist: {{.PolicyTitle}} ({{.PolicyId}})

{{end}}{{if .Author}}**Author:** {{.Author}}{{if .AuthorVersion}} (v{{.AuthorVersion}}){{end}}

{{end}}{{range $index, $section := .Sections}}{{if $index}}
---

{{end}}## Assessment Requirement: {{$section.RequirementId}}

{{if eq (len $section.Items) 0}}- [ ] No evaluation methods defined
{{else}}{{range $section.Items}}- [ ] {{if .MethodDescription}}{{.MethodDescription}}{{else if .MethodType}}{{.MethodType}}{{else}}Evaluation Method{{end}}{{if .Frequency}} ({{.Frequency}}){{end}}{{if .PlanId}} [Plan: {{.PlanId}}]{{end}}
{{if .EvidenceRequirements}}    > **Evidence Required:** {{.EvidenceRequirements}}
{{end}}{{end}}{{end}}{{end}}`
