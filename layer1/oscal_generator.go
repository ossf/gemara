package layer1

import (
	"fmt"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

type generateOpts struct {
	version       string
	imports       map[string]string
	canonicalHref string
}

func (g *generateOpts) complete(doc GuidanceDocument) {
	if g.version == "" {
		g.version = doc.Metadata.Version
	}
	if g.imports == nil {
		g.imports = make(map[string]string)
		for _, mappingRef := range doc.Metadata.MappingReferences {
			g.imports[mappingRef.Id] = mappingRef.Url
		}
	}
}

// GenerateOption defines an option to tune the behavior of the OSCAL
// generation methods for Layer 1.
type GenerateOption func(opts *generateOpts)

// WithVersion is a GenerateOption that sets the version of the OSCAL Document. If set,
// this will be used instead of the version in GuidanceDocument.
func WithVersion(version string) GenerateOption {
	return func(opts *generateOpts) {
		opts.version = version
	}
}

// WithOSCALImports is a GenerateOption that provides the `href` to guidance document mappings in OSCAL
// by mapping unique identifier. If unset, the mapping URL of the guidance document will be used.
func WithOSCALImports(imports map[string]string) GenerateOption {
	return func(opts *generateOpts) {
		opts.imports = imports
	}
}

// WithCanonicalHrefFormat is a GenerateOption that provides an `href` format string
// for the canonical version of the guidance document. If set, this will be added as a
// link in the metadata with the rel="canonical" attribute. Ex - https://myguidance.org/versions/%s
func WithCanonicalHrefFormat(canonicalHref string) GenerateOption {
	return func(opts *generateOpts) {
		opts.canonicalHref = canonicalHref
	}
}

// ToOSCALProfile creates an OSCAL Profile from the imported and local guidelines from
// Layer 1 Guidance Document with a given location to the OSCAL Catalog for the guidance document.
func (g *GuidanceDocument) ToOSCALProfile(guidanceDocHref string, opts ...GenerateOption) (oscal.Profile, error) {
	options := generateOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	options.complete(*g)

	metadata, err := createMetadata(g, options)
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
	// `ToOSCALCatalog` would need to be used to create an OSCAL Catalog for the document.
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

// ToOSCALCatalog creates an OSCAL Catalog from the locally defined guidelines in a given
// Layer 1 Guidance Document.
func (g *GuidanceDocument) ToOSCALCatalog(opts ...GenerateOption) (oscal.Catalog, error) {
	// Return early for empty documents
	if len(g.Categories) == 0 {
		return oscal.Catalog{}, fmt.Errorf("document %s does not have defined guidance categories", g.Metadata.Id)
	}

	options := generateOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	options.complete(*g)

	metadata, err := createMetadata(g, options)
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
		groups = append(groups, g.createControlGroup(category, resourcesMap))
	}

	catalog := oscal.Catalog{
		UUID:       uuid.NewUUID(),
		Metadata:   metadata,
		Groups:     oscalUtils.NilIfEmpty(groups),
		BackMatter: backmatter,
	}
	return catalog, nil
}

func createMetadata(guidance *GuidanceDocument, opts generateOpts) (oscal.Metadata, error) {
	fallbackTime := time.Now()
	metadata := oscal.Metadata{
		Title:        guidance.Metadata.Title,
		OscalVersion: oscal.Version,
		Version:      opts.version,
		Published:    oscalUtils.GetTime(guidance.Metadata.PublicationDate),
		LastModified: oscalUtils.GetTimeWithFallback(guidance.Metadata.LastModified, fallbackTime),
	}

	if opts.canonicalHref != "" {
		metadata.Links = &[]oscal.Link{
			{
				Href: fmt.Sprintf(opts.canonicalHref, opts.version),
				Rel:  "canonical",
			},
		}
	}

	authorRole := oscal.Role{
		ID:          "author",
		Description: "Author and owner of the document",
		Title:       "Author",
	}

	author := oscal.Party{
		UUID: uuid.NewUUID(),
		Type: "person",
		Name: guidance.Metadata.Author,
	}

	responsibleParty := oscal.ResponsibleParty{
		PartyUuids: []string{author.UUID},
		RoleId:     authorRole.ID,
	}

	metadata.Parties = &[]oscal.Party{author}
	metadata.Roles = &[]oscal.Role{authorRole}
	metadata.ResponsibleParties = &[]oscal.ResponsibleParty{responsibleParty}
	return metadata, nil
}

func (g *GuidanceDocument) createControlGroup(category Category, resourcesMap map[string]string) oscal.Group {
	group := oscal.Group{
		Class: "category",
		ID:    category.Id,
		Title: category.Title,
	}

	controlMap := make(map[string]oscal.Control)
	for _, guideline := range category.Guidelines {
		control, parent := g.guidelineToControl(guideline, resourcesMap)

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

func (g *GuidanceDocument) guidelineToControl(guideline Guideline, resourcesMap map[string]string) (oscal.Control, string) {
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

func mappingToLinks(mappings []Mapping, resourcesMap map[string]string) []oscal.Link {
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

func mappingToBackMatter(resourceRefs []MappingReference) *oscal.BackMatter {
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
					"%s. *%s*. %s",
					ref.Issuer,
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
