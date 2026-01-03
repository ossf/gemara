package oscal

import (
	"fmt"
	"sort"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/ossf/gemara"
	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

// FromGuidance creates both an OSCAL Catalog and Profile from a Guidance Document.
// The catalog includes only the locally defined guidelines (categories), not imported ones.
// The profile includes imports for both external guidelines and the local catalog.
func FromGuidance(g *gemara.GuidanceDocument, guidanceDocHref string, opts ...GenerateOption) (oscal.Catalog, oscal.Profile, error) {
	// The guidanceDocHref parameter specifies the location where the OSCAL Catalog
	// will be saved, used to create the import reference in the Profile. This must
	// be a relative or absolute URI that accurately reflects where the catalog
	// file will be located relative to the profile.
	if guidanceDocHref == "" {
		return oscal.Catalog{}, oscal.Profile{}, fmt.Errorf("guidanceDocHref is required to create a valid Profile import reference")
	}
	options := generateOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	options.completeFromGuidance(*g)

	// Create catalog
	// Return early for empty documents
	if len(g.Families) == 0 {
		return oscal.Catalog{}, oscal.Profile{}, fmt.Errorf("document %s does not have defined families", g.Metadata.Id)
	}

	catalogMetadata, err := createMetadataFromGuidance(g, options)
	if err != nil {
		return oscal.Catalog{}, oscal.Profile{}, fmt.Errorf("error creating catalog metadata: %w", err)
	}

	// Create a resource map for control linking
	resourcesMap := make(map[string]string)
	backmatter := mappingToBackMatter(g.Metadata.MappingReferences)
	if backmatter != nil && backmatter.Resources != nil {
		for _, resource := range *backmatter.Resources {
			if resource.Props != nil && len(*resource.Props) > 0 {
				props := *resource.Props
				id := props[0].Value
				resourcesMap[id] = resource.UUID
			}
		}
	}

	// Group guidelines by family
	guidelinesByFamily := make(map[string][]gemara.Guideline)
	for _, guideline := range g.Guidelines {
		// Skip guidelines that extend external controls - these belong only in the profile as alterations
		if guideline.Extends != nil && guideline.Extends.ReferenceId != "" {
			continue
		}
		guidelinesByFamily[guideline.Family] = append(guidelinesByFamily[guideline.Family], guideline)
	}

	var groups []oscal.Group
	for _, family := range g.Families {
		guidelines := guidelinesByFamily[family.Id]
		if len(guidelines) > 0 {
			groups = append(groups, createControlGroup(g, family, guidelines, resourcesMap))
		}
	}

	catalog := oscal.Catalog{
		UUID:       uuid.NewUUID(),
		Metadata:   catalogMetadata,
		Groups:     oscalUtils.NilIfEmpty(groups),
		BackMatter: backmatter,
	}

	profileMetadata, err := createMetadataFromGuidance(g, options)
	if err != nil {
		return oscal.Catalog{}, oscal.Profile{}, fmt.Errorf("error creating profile metadata: %w", err)
	}

	importMap := make(map[string]oscal.Import)
	for mappingId, mappingRef := range options.imports {
		importMap[mappingId] = oscal.Import{Href: mappingRef}
	}

	alterationMap := processExternalControls(g.Guidelines, importMap, options.imports, g.Metadata.Id)

	var imports []oscal.Import
	for _, imp := range importMap {
		if imp.IncludeControls != nil || imp.IncludeAll != nil {
			imports = append(imports, imp)
		}
	}

	// Add an import for each control defined locally in the Guidance Document
	// The catalog is created by FromGuidance and referenced here.
	localImport := oscal.Import{
		Href:       guidanceDocHref,
		IncludeAll: &oscal.IncludeAll{},
	}
	imports = append(imports, localImport)

	modify := buildModifySection(alterationMap)

	profile := oscal.Profile{
		UUID:     uuid.NewUUID(),
		Imports:  imports,
		Metadata: profileMetadata,
		Modify:   modify,
	}

	return catalog, profile, nil
}

func createControlGroup(g *gemara.GuidanceDocument, family gemara.Family, guidelines []gemara.Guideline, resourcesMap map[string]string) oscal.Group {
	group := oscal.Group{
		Class: "family",
		ID:    family.Id,
		Title: family.Title,
	}

	controlMap := make(map[string]oscal.Control)
	parentChildMap := make(map[string][]string)
	childControlIds := make(map[string]struct{})

	// Create all controls and track parent-child relationships
	for _, guideline := range guidelines {
		control, parent := guidelineToControl(g, guideline, resourcesMap)
		controlMap[control.ID] = control

		if parent != "" {
			parentChildMap[parent] = append(parentChildMap[parent], control.ID)
			childControlIds[control.ID] = struct{}{}
		}
	}

	// Link children to their parents using a queue strategy
	queue := make([]string, 0, len(parentChildMap))
	processed := make(map[string]bool)

	for parentId := range parentChildMap {
		queue = append(queue, parentId)
	}

	for len(queue) > 0 {
		parentId := queue[0]
		queue = queue[1:]

		if processed[parentId] {
			continue
		}

		parentControl, exists := controlMap[parentId]
		if !exists {
			// Drop orphaned controls
			processed[parentId] = true
			continue
		}

		childIds := parentChildMap[parentId]
		allChildrenReady := true

		// Check if any children are themselves parents that haven't been processed yet
		for _, childId := range childIds {
			if _, isParent := parentChildMap[childId]; isParent && !processed[childId] {
				allChildrenReady = false
				break
			}
		}

		if !allChildrenReady {
			// Add back to queue to process later
			queue = append(queue, parentId)
			continue
		}

		if parentControl.Controls == nil {
			parentControl.Controls = &[]oscal.Control{}
		}
		children := make([]oscal.Control, 0, len(childIds))
		for _, childId := range childIds {
			if childControl, exists := controlMap[childId]; exists {
				children = append(children, childControl)
			}
		}
		parentControl.Controls = &children
		controlMap[parentId] = parentControl
		processed[parentId] = true
	}

	controls := make([]oscal.Control, 0, len(controlMap))
	for id, control := range controlMap {
		if _, isChild := childControlIds[id]; !isChild {
			controls = append(controls, control)
		}
	}

	group.Controls = oscalUtils.NilIfEmpty(controls)
	return group
}

// guidelineToParts converts a guideline to OSCAL parts that can be added to an existing control.
// This is used when a guideline extends an existing control via the profile's modify mechanism.
// If guidelineId is provided, parts use the naming convention: {controlId}_{guidelineId}_{partType}
// If guidelineId is empty, parts use the standard naming convention: {controlId}_{partType}
func guidelineToParts(guideline gemara.Guideline, controlId string, guidelineId string) []oscal.Part {
	var parts []oscal.Part

	// Determine the part ID prefix based on whether this is an alteration or a full control
	var prefix string
	if guidelineId != "" {
		// For alterations: {controlId}_{guidelineId}
		normalizedGuidelineId := oscalUtils.NormalizeControl(guidelineId, false)
		prefix = fmt.Sprintf("%s_%s", controlId, normalizedGuidelineId)
	} else {
		// For full controls: {controlId}
		prefix = controlId
	}

	// Add overview part if objective exists
	if guideline.Objective != "" {
		parts = append(parts, oscal.Part{
			Name:  "overview",
			ID:    fmt.Sprintf("%s_ovw", prefix),
			Prose: guideline.Objective,
		})
	}

	// Build statement parts
	var statementParts []oscal.Part
	var objectiveParts []oscal.Part
	for _, statement := range guideline.Statements {
		partId := oscalUtils.NormalizeControl(statement.Id, true)
		statementId := fmt.Sprintf("%s_smt.%s", prefix, partId)

		statementParts = append(statementParts, oscal.Part{
			Name:  "item",
			ID:    statementId,
			Prose: statement.Text,
			Title: statement.Title,
		})

		if len(statement.Recommendations) > 0 {
			objectiveId := fmt.Sprintf("%s_obj.%s", prefix, partId)
			objectiveParts = append(objectiveParts, oscal.Part{
				Name:  "assessment-objective",
				ID:    objectiveId,
				Prose: strings.Join(statement.Recommendations, " "),
				Links: &[]oscal.Link{
					{
						Href: fmt.Sprintf("#%s", statementId),
						Rel:  "assessment-for",
					},
				},
			})
		}
	}

	if len(statementParts) > 0 {
		parts = append(parts, oscal.Part{
			Name:  "statement",
			ID:    fmt.Sprintf("%s_smt", prefix),
			Parts: &statementParts,
		})
	}

	if len(guideline.Recommendations) > 0 || len(objectiveParts) > 0 {
		objectivePart := oscal.Part{
			Name: "assessment-objective",
			ID:   fmt.Sprintf("%s_obj", prefix),
		}
		if len(guideline.Recommendations) > 0 {
			objectivePart.Prose = strings.Join(guideline.Recommendations, " ")
			objectivePart.Links = &[]oscal.Link{
				{
					Href: fmt.Sprintf("#%s_smt", prefix),
					Rel:  "assessment-for",
				},
			}
		}
		if len(objectiveParts) > 0 {
			objectivePart.Parts = &objectiveParts
		}
		parts = append(parts, objectivePart)
	}

	return parts
}

func guidelineToControl(g *gemara.GuidanceDocument, guideline gemara.Guideline, resourcesMap map[string]string) (oscal.Control, string) {
	controlId := oscalUtils.NormalizeControl(guideline.Id, false)

	control := oscal.Control{
		ID:    controlId,
		Title: guideline.Title,
		Class: g.Metadata.Id,
	}

	var links []oscal.Link
	for _, also := range guideline.SeeAlso {
		relatedLink := oscal.Link{
			Href: fmt.Sprintf("#%s", oscalUtils.NormalizeControl(also, false)),
			Rel:  "related",
		}
		links = append(links, relatedLink)
	}

	guidanceLinks := mappingToLinks(guideline.GuidelineMappings, resourcesMap)
	principleLinks := mappingToLinks(guideline.PrincipleMappings, resourcesMap)
	links = append(links, guidanceLinks...)
	links = append(links, principleLinks...)
	control.Links = oscalUtils.NilIfEmpty(links)

	parts := guidelineToParts(guideline, controlId, "")
	control.Parts = completeParts(parts, controlId)

	parentId := ""
	// Only process local controls (no reference id)
	if guideline.Extends != nil && guideline.Extends.ReferenceId == "" {
		parentId = guideline.Extends.EntryId
	}
	return control, oscalUtils.NormalizeControl(parentId, false)
}

// processAlteration creates or updates an alteration for a control and merges guideline parts into it.
func processAlteration(alterationMap map[string]*oscal.Alteration, normalizedId string, guideline gemara.Guideline, frameworkPrefix string) {
	alteration, exists := alterationMap[normalizedId]
	if !exists {
		alteration = &oscal.Alteration{
			ControlId: normalizedId,
		}
		alterationMap[normalizedId] = alteration
	}

	// Generate parts from guideline using framework prefix for consistent naming
	parts := guidelineToParts(guideline, normalizedId, frameworkPrefix)
	if len(parts) == 0 {
		return
	}

	if alteration.Adds == nil {
		alteration.Adds = &[]oscal.Addition{}
	}

	// Merge parts into existing addition if one exists, otherwise create new
	if len(*alteration.Adds) > 0 {
		firstAddition := &(*alteration.Adds)[0]
		if firstAddition.Parts == nil {
			firstAddition.Parts = &parts
		} else {
			*firstAddition.Parts = append(*firstAddition.Parts, parts...)
		}
	} else {
		addition := oscal.Addition{
			Parts: &parts,
		}
		*alteration.Adds = append(*alteration.Adds, addition)
	}
}

// processExternalControls processes guidelines that extend external controls.
func processExternalControls(guidelines []gemara.Guideline, importMap map[string]oscal.Import, imports map[string]string, frameworkPrefix string) map[string]*oscal.Alteration {
	alterationMap := make(map[string]*oscal.Alteration)

	for _, guideline := range guidelines {
		// Do not process guidelines that extend local controls.
		// This is handled in catalogs.
		if guideline.Extends == nil || guideline.Extends.ReferenceId == "" || guideline.Extends.EntryId == "" {
			continue
		}

		href, found := imports[guideline.Extends.ReferenceId]
		if !found || href == "" {
			continue
		}

		imp := getOrCreateImport(importMap, guideline.Extends.ReferenceId, href)
		normalizedId := oscalUtils.NormalizeControl(guideline.Extends.EntryId, false)
		imp.IncludeControls = mergeControlIds(imp.IncludeControls, normalizedId)
		importMap[guideline.Extends.ReferenceId] = imp

		processAlteration(alterationMap, normalizedId, guideline, frameworkPrefix)
	}

	return alterationMap
}

// getOrCreateImport retrieves an existing import or creates a new one.
func getOrCreateImport(importMap map[string]oscal.Import, referenceId, href string) oscal.Import {
	if imp, exists := importMap[referenceId]; exists {
		return imp
	}
	return oscal.Import{Href: href}
}

// mergeControlIds merges a new control ID with existing control selectors, returning
// a single selector containing all unique control IDs sorted.
func mergeControlIds(existingSelectors *[]oscal.SelectControlById, newControlId string) *[]oscal.SelectControlById {
	allControlIds := make(map[string]struct{})

	if existingSelectors != nil {
		for _, selector := range *existingSelectors {
			if selector.WithIds != nil {
				for _, id := range *selector.WithIds {
					allControlIds[id] = struct{}{}
				}
			}
		}
	}

	allControlIds[newControlId] = struct{}{}

	mergedIds := make([]string, 0, len(allControlIds))
	for id := range allControlIds {
		mergedIds = append(mergedIds, id)
	}
	sort.Strings(mergedIds)

	selector := oscal.SelectControlById{WithIds: &mergedIds}
	return &[]oscal.SelectControlById{selector}
}

// buildModifySection creates a Modify section from alterations if any exist.
func buildModifySection(alterationMap map[string]*oscal.Alteration) *oscal.Modify {
	if len(alterationMap) == 0 {
		return nil
	}

	alterations := make([]oscal.Alteration, 0, len(alterationMap))
	for _, alteration := range alterationMap {
		alterations = append(alterations, *alteration)
	}

	return &oscal.Modify{
		Alters: &alterations,
	}
}

// completeParts ensures that statement and assessment-objective parts always exist
// as required by OSCAL, creating empty defaults if needed.
func completeParts(parts []oscal.Part, controlId string) *[]oscal.Part {
	var statementPart oscal.Part
	var assessmentObjectivePart oscal.Part
	var otherParts []oscal.Part

	for _, part := range parts {
		switch part.Name {
		case "statement":
			statementPart = part
		case "assessment-objective":
			assessmentObjectivePart = part
		default:
			otherParts = append(otherParts, part)
		}
	}

	if statementPart.ID == "" {
		statementPart = oscal.Part{
			Name: "statement",
			ID:   fmt.Sprintf("%s_smt", controlId),
		}
	}

	if assessmentObjectivePart.ID == "" {
		assessmentObjectivePart = oscal.Part{
			Name: "assessment-objective",
			ID:   fmt.Sprintf("%s_obj", controlId),
		}
	}

	finalParts := []oscal.Part{statementPart, assessmentObjectivePart}
	finalParts = append(finalParts, otherParts...)
	return &finalParts
}

func mappingToLinks(mappings []gemara.MultiMapping, resourcesMap map[string]string) []oscal.Link {
	links := make([]oscal.Link, 0, len(mappings))
	for _, mapping := range mappings {
		ref, found := resourcesMap[mapping.ReferenceId]
		if !found {
			continue
		}
		externalLink := oscal.Link{
			Href: fmt.Sprintf("#%s", ref),
			Rel:  "reference",
		}
		links = append(links, externalLink)
	}
	return links
}

func mappingToBackMatter(resourceRefs []gemara.MappingReference) *oscal.BackMatter {
	var resources []oscal.Resource
	for _, ref := range resourceRefs {
		resource := oscal.Resource{
			UUID:        uuid.NewUUID(),
			Title:       ref.Title,
			Description: ref.Description,
			Props: &[]oscal.Property{
				{
					Name:  "id",
					Value: ref.Id,
					Ns:    oscalUtils.GemaraNamespace,
				},
			},
			Rlinks: &[]oscal.ResourceLink{
				{
					Href: ref.Url,
				},
			},
			Citation: &oscal.Citation{
				Text: fmt.Sprintf(
					"*%s*. %s",
					ref.Title,
					ref.Url),
			},
		}
		resources = append(resources, resource)
	}

	if len(resources) == 0 {
		return nil
	}

	backmatter := oscal.BackMatter{
		Resources: &resources,
	}
	return &backmatter
}
