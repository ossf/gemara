package checklist

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ossf/gemara"
)

// Item represents a single checklist item.
type Item struct {
	// RequirementId is the requirement ID (e.g., "OSPS-AC-01.01")
	RequirementId string
	// ProcedureName is the human-readable name of the procedure to execute.
	EvaluatorName string
	// Description provides additional context or a summary about the procedure.
	Description string
	// Documentation is the documentation URL
	Documentation string
	// IsAdditionalProcedure indicates if this is an additional procedure
	IsAdditionalEvaluator bool
}

// RequirementSection organizes checklist items by assessment requirement.
type RequirementSection struct {
	// RequirementId is the reference to the unique assessment requirement identifier.
	RequirementId string
	// RequirementRecommendation is the reference to the recommendation or modified recommendation
	// for assessment requirement implementation.
	RequirementRecommendation string
	// Items are the checklist items for this control
	Items []Item
}

type ControlSection struct {
	ControlRef   string
	Requirements []RequirementSection
}

// Checklist represents the structured checklist data.
type Checklist struct {
	// PolicyId identifies the evaluation plan.
	PolicyId string
	// Author is the name of the plan author.
	Author string
	// AuthorVersion is the version of the authoring tool or system.
	AuthorVersion string
	// Sections are the control sections
	Sections []ControlSection
}

// ToChecklist converts an EvaluationPlan into a structured Checklist.
// Catalogs are loaded on-demand from metadata.mapping-references when needed.
func ToChecklist(policy gemara.Policy) (Checklist, error) {
	checklist := Checklist{}

	if policy.Metadata.Id != "" {
		checklist.PolicyId = policy.Metadata.Id
	}
	if policy.Metadata.Author.Name != "" {
		checklist.Author = policy.Metadata.Author.Name
		checklist.AuthorVersion = policy.Metadata.Author.Version
	}

	// Cache for catalogs loaded on-demand
	catalogCache := make(map[string]*gemara.Catalog)

	for _, controlRef := range policy.ControlReferences {
		items, err := buildChecklistItems(policy.ImplementationPlan.Evaluators, controlRef.AssessmentRequirementModifications, controlRef.ReferenceId, policy.Metadata.MappingReferences, catalogCache)
		if err != nil {
			return Checklist{}, fmt.Errorf("failed to build checklist items for control reference %q: %w", controlRef.ReferenceId, err)
		}

		section := ControlSection{
			ControlRef:   controlRef.ReferenceId,
			Requirements: items,
		}
		checklist.Sections = append(checklist.Sections, section)
	}

	return checklist, nil
}

// ToMarkdownChecklist converts an evaluation plan into a checklist.
// Generates a pre-execution checklist showing what needs to be checked.
func ToMarkdownChecklist(checklist Checklist) (string, error) {
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

// findAssessmentRequirement finds an AssessmentRequirement by ID in a catalog.
func findAssessmentRequirement(catalog *gemara.Catalog, requirementID string) *gemara.AssessmentRequirement {
	if catalog == nil {
		return nil
	}

	for _, family := range catalog.ControlFamilies {
		for i := range family.Controls {
			control := &family.Controls[i]
			for j := range control.AssessmentRequirements {
				requirement := &control.AssessmentRequirements[j]
				if requirement.Id == requirementID {
					return requirement
				}
			}
		}
	}

	return nil
}

// loadCatalogOnDemand loads a catalog from the mapping references if not already cached.
func loadCatalogOnDemand(refId string, mappingReferences []gemara.MappingReference, cache map[string]*gemara.Catalog) (*gemara.Catalog, error) {
	// Check cache first
	if catalog, ok := cache[refId]; ok {
		return catalog, nil
	}

	// Find the mapping reference
	var mappingRef *gemara.MappingReference
	for i := range mappingReferences {
		if mappingReferences[i].Id == refId {
			mappingRef = &mappingReferences[i]
			break
		}
	}

	if mappingRef == nil || mappingRef.Url == "" {
		// No mapping reference found or no URL - return nil (catalog not available)
		return nil, nil
	}

	// Ensure the URL has a file:// scheme and is absolute
	catalogUrl := mappingRef.Url
	if !strings.HasPrefix(catalogUrl, "file://") && !strings.HasPrefix(catalogUrl, "https://") {
		// Relative path - convert to absolute
		absPath, err := filepath.Abs(catalogUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve catalog path %q: %w", catalogUrl, err)
		}
		catalogUrl = "file://" + absPath
	} else if strings.HasPrefix(catalogUrl, "file://") {
		// Already has file:// scheme, but may be relative
		filePath := strings.TrimPrefix(catalogUrl, "file://")
		if !filepath.IsAbs(filePath) {
			absPath, err := filepath.Abs(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to resolve catalog path %q: %w", filePath, err)
			}
			catalogUrl = "file://" + absPath
		}
	}

	// Load the catalog
	catalog := &gemara.Catalog{}
	if err := catalog.LoadFile(catalogUrl); err != nil {
		return nil, fmt.Errorf("failed to load catalog from %q: %w", catalogUrl, err)
	}

	// Cache it
	cache[refId] = catalog
	return catalog, nil
}

// buildChecklistItems converts an Assessment Requirement into checklist items.
func buildChecklistItems(evaluators []gemara.Actor, modifiers []gemara.AssessmentRequirementModifier, refId string, mappingReferences []gemara.MappingReference, catalogCache map[string]*gemara.Catalog) ([]RequirementSection, error) {
	var items []RequirementSection
	assessmentNum := 1

	evaluatorsMap := make(map[string]gemara.Actor)
	for _, evaluator := range evaluators {
		evaluatorsMap[evaluator.Id] = evaluator
	}

	for _, assessment := range modifiers {
		requirementId := assessment.TargetId
		if requirementId == "" {
			requirementId = fmt.Sprintf("Assessment %d", assessmentNum)
		}
		assessmentNum++

		// Check if recommendation is overridden
		recommendationOverridden := assessment.Overrides != nil && assessment.Overrides.Recommendation != ""

		// Add description item for the AssessmentRequirement
		// Use override text if available, otherwise load from Layer 2 if recommendation is not overridden
		var descriptionText string
		if assessment.Overrides != nil && assessment.Overrides.Text != "" {
			descriptionText = assessment.Overrides.Text
		} else if !recommendationOverridden {
			catalog, err := loadCatalogOnDemand(refId, mappingReferences, catalogCache)
			if err != nil {
				return nil, fmt.Errorf("failed to load catalog for requirement %q: %w", requirementId, err)
			}
			if catalog != nil {
				layer2Requirement := findAssessmentRequirement(catalog, requirementId)
				if layer2Requirement != nil && layer2Requirement.Text != "" {
					descriptionText = layer2Requirement.Text
				}
			}
		}

		requirementSection := RequirementSection{
			RequirementId:             requirementId,
			RequirementRecommendation: descriptionText,
			Items:                     []Item{},
		}

		if assessment.Extensions == nil || len(assessment.Extensions.RequiredEvaluators) == 0 {
			return nil, fmt.Errorf("assessment %q has no evaluators", requirementId)
		}

		for i, reqEval := range assessment.Extensions.RequiredEvaluators {
			evaluator := evaluatorsMap[reqEval]

			item := Item{
				RequirementId:         requirementId,
				EvaluatorName:         evaluator.Name,
				Description:           evaluator.Description,
				Documentation:         evaluator.Uri,
				IsAdditionalEvaluator: i > 0,
			}
			requirementSection.Items = append(requirementSection.Items, item)
		}

		items = append(items, requirementSection)
	}

	return items, nil
}
