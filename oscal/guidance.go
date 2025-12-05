package oscal

import (
	"fmt"
	"strings"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/ossf/gemara"
	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

// ProfileFromGuidanceDocument creates an OSCAL Profile from the imported and local guidelines from
// Layer 1 Guidance Document with a given location to the OSCAL Catalog for the guidance document.
func ProfileFromGuidanceDocument(g *gemara.GuidanceDocument, guidanceDocHref string, opts ...GenerateOption) (oscal.Profile, error) {
	options := generateOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	options.completeFromGuidance(*g)

	metadata, err := createMetadataFromGuidance(g, options)
	if err != nil {
		return oscal.Profile{}, fmt.Errorf("error creating profile metadata: %w", err)
	}

	importMap := make(map[string]oscal.Import)
	for mappingId, mappingRef := range options.imports {
		importMap[mappingId] = oscal.Import{Href: mappingRef}
	}

	for _, mapping := range g.ImportedGuidelines {
		imp, ok := importMap[mapping.ReferenceId]
		if !ok {
			continue
		}

		withIds := make([]string, 0, len(mapping.Entries))
		for _, entry := range mapping.Entries {
			withIds = append(withIds, oscalUtils.NormalizeControl(entry.ReferenceId, false))
		}

		selector := oscal.SelectControlById{WithIds: &withIds}
		imp.IncludeControls = &[]oscal.SelectControlById{selector}
		importMap[mapping.ReferenceId] = imp
	}

	var imports []oscal.Import
	for _, imp := range importMap {
		if imp.IncludeControls != nil {
			imports = append(imports, imp)
		}
	}

	// Add an import for each control defined locally in the Layer 1 Guidance Document
	// `CatalogFromGuidanceDocument` would need to be used to create an OSCAL Catalog for the document.
	localImport := oscal.Import{
		Href:       guidanceDocHref,
		IncludeAll: &oscal.IncludeAll{},
	}
	imports = append(imports, localImport)

	profile := oscal.Profile{
		UUID:     uuid.NewUUID(),
		Imports:  imports,
		Metadata: metadata,
	}
	return profile, nil
}

// CatalogFromGuidanceDocument creates an OSCAL Catalog from the locally defined guidelines in a given
// Layer 1 Guidance Document.
func CatalogFromGuidanceDocument(g *gemara.GuidanceDocument, opts ...GenerateOption) (oscal.Catalog, error) {
	// Return early for empty documents
	if len(g.Categories) == 0 {
		return oscal.Catalog{}, fmt.Errorf("document %s does not have defined guidance categories", g.Metadata.Id)
	}

	options := generateOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	options.completeFromGuidance(*g)

	metadata, err := createMetadataFromGuidance(g, options)
	if err != nil {
		return oscal.Catalog{}, fmt.Errorf("error creating catalog metadata: %w", err)
	}

	// Create a resource map for control linking
	resourcesMap := make(map[string]string)
	backmatter := mappingToBackMatter(g.Metadata.MappingReferences)
	if backmatter != nil {
		for _, resource := range *backmatter.Resources {
			// Extract the id from the props
			props := *resource.Props
			id := props[0].Value
			resourcesMap[id] = resource.UUID
		}
	}

	var groups []oscal.Group
	for _, category := range g.Categories {
		groups = append(groups, createControlGroup(g, category, resourcesMap))
	}

	catalog := oscal.Catalog{
		UUID:       uuid.NewUUID(),
		Metadata:   metadata,
		Groups:     oscalUtils.NilIfEmpty(groups),
		BackMatter: backmatter,
	}
	return catalog, nil
}

func createControlGroup(g *gemara.GuidanceDocument, category gemara.Category, resourcesMap map[string]string) oscal.Group {
	group := oscal.Group{
		Class: "category",
		ID:    category.Id,
		Title: category.Title,
	}

	controlMap := make(map[string]oscal.Control)
	for _, guideline := range category.Guidelines {
		control, parent := guidelineToControl(g, guideline, resourcesMap)

		if parent == "" {
			controlMap[control.ID] = control
		} else {
			parentControl := controlMap[parent]
			if parentControl.Controls == nil {
				parentControl.Controls = &[]oscal.Control{}
			}
			*parentControl.Controls = append(*parentControl.Controls, control)
			controlMap[parent] = parentControl
		}
	}

	controls := make([]oscal.Control, 0, len(controlMap))
	for _, control := range controlMap {
		controls = append(controls, control)
	}

	group.Controls = oscalUtils.NilIfEmpty(controls)
	return group
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

	// Top-level statements are required for controls per OSCAL guidance
	smtPart := oscal.Part{
		Name: "statement",
		ID:   fmt.Sprintf("%s_smt", controlId),
	}

	objPart := oscal.Part{
		Name: "assessment-objective",
		ID:   fmt.Sprintf("%s_obj", controlId),
	}

	if len(guideline.Recommendations) > 0 {
		objPart.Prose = strings.Join(guideline.Recommendations, " ")
		objPart.Links = &[]oscal.Link{
			{
				Href: fmt.Sprintf("#%s_smt", controlId),
				Rel:  "assessment-for",
			},
		}
	}

	var smtParts []oscal.Part
	var objParts []oscal.Part
	for _, part := range guideline.GuidelineParts {
		partId := oscalUtils.NormalizeControl(part.Id, true)
		smtID := fmt.Sprintf("%s_smt.%s", controlId, partId)
		itemSubSmt := oscal.Part{
			Name:  "item",
			ID:    smtID,
			Prose: part.Text,
			Title: part.Title,
		}
		smtParts = append(smtParts, itemSubSmt)

		if len(part.Recommendations) > 0 {
			objSubPart := oscal.Part{
				Name:  "assessment-objective",
				ID:    fmt.Sprintf("%s_obj.%s", controlId, partId),
				Prose: strings.Join(part.Recommendations, " "),
				Links: &[]oscal.Link{
					{
						Href: fmt.Sprintf("#%s", smtID),
						Rel:  "assessment-for",
					},
				},
			}
			objParts = append(objParts, objSubPart)
		}
	}

	// Ensure the parts are set to nil if nothing was added for
	// schema compliance.
	smtPart.Parts = oscalUtils.NilIfEmpty(smtParts)
	objPart.Parts = oscalUtils.NilIfEmpty(objParts)
	control.Parts = &[]oscal.Part{smtPart, objPart}

	if guideline.Objective != "" {
		overviewPart := oscal.Part{
			Name:  "overview",
			ID:    fmt.Sprintf("%s_ovw", controlId),
			Prose: guideline.Objective,
		}
		*control.Parts = append(*control.Parts, overviewPart)
	}

	return control, oscalUtils.NormalizeControl(guideline.BaseGuidelineID, false)
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
