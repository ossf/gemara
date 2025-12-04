package oscal

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/ossf/gemara"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

const defaultVersion = "0.0.1"

type generateOpts struct {
	version       string
	imports       map[string]string
	canonicalHref string
}

func (g *generateOpts) complete(doc gemara.Guidance) {
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
// generation functions for Layer 1.
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
// link in the mapping.cue with the rel="canonical" attribute. Ex - https://myguidance.org/versions/%s
func WithCanonicalHrefFormat(canonicalHref string) GenerateOption {
	return func(opts *generateOpts) {
		opts.canonicalHref = canonicalHref
	}
}

// ProfileFromGuidanceDocument creates an OSCAL Profile from the imported and local guidelines from
// Layer 1 Guidance Document with a given location to the OSCAL Catalog for the guidance document.
func ProfileFromGuidanceDocument(g *gemara.Guidance, guidanceDocHref string, opts ...GenerateOption) (oscal.Profile, error) {
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

	// Process extends mappings from guidelines as imports and "add" modifications
	alterationMap := make(map[string]*oscal.Alteration)

	for _, category := range g.Categories {
		for _, guideline := range category.Guidelines {
			refId := guideline.Extends.ReferenceId
			if refId != "" && guideline.Extends.EntryId != "" {
				href, found := options.imports[refId]
				if found && href != "" {
					imp, exists := importMap[refId]
					if !exists {
						imp = oscal.Import{Href: href}
					}

					normalizedId := oscalUtils.NormalizeControl(guideline.Extends.EntryId, false)
					withIds := []string{normalizedId}

					// Merge with existing IncludeControls if any
					if imp.IncludeControls == nil {
						imp.IncludeControls = &[]oscal.SelectControlById{}
					}

					// Check if we already have a selector for this set of controls
					// If not, create a new one and merge all control IDs
					allControlIds := make(map[string]bool)
					for _, selector := range *imp.IncludeControls {
						if selector.WithIds != nil {
							for _, id := range *selector.WithIds {
								allControlIds[id] = true
							}
						}
					}
					for _, id := range withIds {
						allControlIds[id] = true
					}

					// Create a single selector with all unique control IDs
					mergedIds := make([]string, 0, len(allControlIds))
					for id := range allControlIds {
						mergedIds = append(mergedIds, id)
					}
					sort.Strings(mergedIds)
					selector := oscal.SelectControlById{WithIds: &mergedIds}
					imp.IncludeControls = &[]oscal.SelectControlById{selector}
					importMap[refId] = imp

					// Create or update alteration for this control
					alteration, exists := alterationMap[normalizedId]
					if !exists {
						alteration = &oscal.Alteration{
							ControlId: normalizedId,
						}
						alterationMap[normalizedId] = alteration
					}

					// Use document ID as prefix for framework-level extensions (e.g., 800-161 extending 800-53)
					// This ensures all enhancements from the same framework use consistent naming
					frameworkPrefix := g.Metadata.Id
					parts := guidelineToParts(guideline, normalizedId, frameworkPrefix)
					if len(parts) > 0 {
						if alteration.Adds == nil {
							alteration.Adds = &[]oscal.Addition{}
						}
						// Merge parts into existing addition if one exists, otherwise create new
						if len(*alteration.Adds) > 0 {
							// Merge parts into the first addition
							existingParts := (*alteration.Adds)[0].Parts
							if existingParts == nil {
								(*alteration.Adds)[0].Parts = &parts
							} else {
								*existingParts = append(*existingParts, parts...)
							}
						} else {
							addition := oscal.Addition{
								Parts: &parts,
							}
							*alteration.Adds = append(*alteration.Adds, addition)
						}
					}
				}
			}
		}
	}

	var imports []oscal.Import
	for _, imp := range importMap {
		if imp.IncludeControls != nil || imp.IncludeAll != nil {
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

	// Build modify section if we have any alterations
	var modify *oscal.Modify
	if len(alterationMap) > 0 {
		alterations := make([]oscal.Alteration, 0, len(alterationMap))
		for _, alteration := range alterationMap {
			alterations = append(alterations, *alteration)
		}
		modify = &oscal.Modify{
			Alters: &alterations,
		}
	}

	profile := oscal.Profile{
		UUID:     uuid.NewUUID(),
		Imports:  imports,
		Metadata: metadata,
		Modify:   modify,
	}
	return profile, nil
}

// CatalogFromGuidanceDocument creates an OSCAL Catalog from the locally defined guidelines in a given
// Layer 1 Guidance Document.
func CatalogFromGuidanceDocument(g *gemara.Guidance, opts ...GenerateOption) (oscal.Catalog, error) {
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

func createMetadata(guidance *gemara.Guidance, opts generateOpts) (oscal.Metadata, error) {
	now := time.Now()
	metadata := oscal.Metadata{
		Title:        guidance.Title,
		OscalVersion: oscal.Version,
		Version:      opts.version,
		Published:    oscalUtils.GetTime(string(guidance.Metadata.Date)),
		LastModified: now,
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
		Name: guidance.Metadata.Author.Name,
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

func createControlGroup(g *gemara.Guidance, category gemara.Category, resourcesMap map[string]string) oscal.Group {
	group := oscal.Group{
		Class: "category",
		ID:    category.Id,
		Title: category.Title,
	}

	controlMap := make(map[string]oscal.Control)
	for _, guideline := range category.Guidelines {
		// Don't process this as a local catalog control
		// It extend another and would need to be
		// resolved through a profile.
		if guideline.Extends.ReferenceId != "" {
			continue
		}
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
		overviewPart := oscal.Part{
			Name:  "overview",
			ID:    fmt.Sprintf("%s_ovw", prefix),
			Prose: guideline.Objective,
		}
		parts = append(parts, overviewPart)
	}

	// Add statement parts if they exist
	var smtParts []oscal.Part
	var objParts []oscal.Part
	for _, part := range guideline.Statements {
		partId := oscalUtils.NormalizeControl(part.Id, true)
		smtID := fmt.Sprintf("%s_smt.%s", prefix, partId)
		itemSubSmt := oscal.Part{
			Name:  "item",
			ID:    smtID,
			Prose: part.Text,
			Title: part.Title,
		}
		smtParts = append(smtParts, itemSubSmt)

		if len(part.Recommendations) > 0 {
			objID := fmt.Sprintf("%s_obj.%s", prefix, partId)
			objSubPart := oscal.Part{
				Name:  "assessment-objective",
				ID:    objID,
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

	// Add statement part with sub-items if any exist
	if len(smtParts) > 0 {
		smtPart := oscal.Part{
			Name:  "statement",
			ID:    fmt.Sprintf("%s_smt", prefix),
			Parts: &smtParts,
		}
		parts = append(parts, smtPart)
	}

	// Add assessment-objective part with sub-items if recommendations exist
	if len(guideline.Recommendations) > 0 || len(objParts) > 0 {
		objPart := oscal.Part{
			Name: "assessment-objective",
			ID:   fmt.Sprintf("%s_obj", prefix),
		}
		if len(guideline.Recommendations) > 0 {
			objPart.Prose = strings.Join(guideline.Recommendations, " ")
			objPart.Links = &[]oscal.Link{
				{
					Href: fmt.Sprintf("#%s_smt", prefix),
					Rel:  "assessment-for",
				},
			}
		}
		if len(objParts) > 0 {
			objPart.Parts = &objParts
		}
		parts = append(parts, objPart)
	}

	return parts
}

func guidelineToControl(g *gemara.Guidance, guideline gemara.Guideline, resourcesMap map[string]string) (oscal.Control, string) {
	controlId := oscalUtils.NormalizeControl(guideline.Id, false)

	control := oscal.Control{
		ID:    controlId,
		Title: guideline.Title,
		Class: g.Metadata.Id,
	}

	var links []oscal.Link
	for _, also := range guideline.SeeAlso {
		relatedLink := oscal.Link{
			Href: fmt.Sprintf("#%s", oscalUtils.NormalizeControl(also.EntryId, false)),
			Rel:  "related",
		}
		links = append(links, relatedLink)
	}

	guidanceLinks := mappingToLinks(guideline.GuidelineMappings, resourcesMap)
	principleLinks := mappingToLinks(guideline.PrincipleMappings, resourcesMap)
	links = append(links, guidanceLinks...)
	links = append(links, principleLinks...)
	control.Links = oscalUtils.NilIfEmpty(links)

	// Use guidelineToParts to generate the parts (empty guidelineId for full controls)
	parts := guidelineToParts(guideline, controlId, "")

	// OSCAL requires statement and assessment-objective parts to always exist
	var statementPart oscal.Part
	var assessmentObjectivePart oscal.Part
	var otherParts []oscal.Part
	hasStatement := false
	hasAssessmentObjective := false

	for _, part := range parts {
		switch part.Name {
		case "statement":
			statementPart = part
			hasStatement = true
		case "assessment-objective":
			assessmentObjectivePart = part
			hasAssessmentObjective = true
		default:
			otherParts = append(otherParts, part)
		}
	}

	// Create default parts if they don't exist
	if !hasStatement {
		statementPart = oscal.Part{
			Name: "statement",
			ID:   fmt.Sprintf("%s_smt", controlId),
		}
	}
	if !hasAssessmentObjective {
		assessmentObjectivePart = oscal.Part{
			Name: "assessment-objective",
			ID:   fmt.Sprintf("%s_obj", controlId),
		}
	}

	finalParts := []oscal.Part{statementPart, assessmentObjectivePart}
	finalParts = append(finalParts, otherParts...)
	control.Parts = &finalParts

	return control, oscalUtils.NormalizeControl(guideline.Extends.EntryId, false)
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

// FromCatalog converts a Catalog to OSCAL Catalog format.
// Parameters:
//   - catalog: The catalog to convert
//   - controlHREF: URL template for linking to controls. Uses format: controlHREF(version, controlID)
//     Example: "https://baseline.openssf.org/versions/%s#%s"
//
// The function automatically:
//   - Uses the catalog's internal version from Metadata.Version
//   - Uses the ControlFamily.Id as the OSCAL group ID
//   - Generates a unique UUID for the catalog
func FromCatalog(catalog *gemara.Catalog, controlHREF string) (oscal.Catalog, error) {
	now := time.Now()

	version := catalog.Metadata.Version
	if catalog.Metadata.Version == "" {
		version = defaultVersion
	}

	oscalCatalog := oscal.Catalog{
		UUID:   uuid.NewUUID(),
		Groups: nil,
		Metadata: oscal.Metadata{
			LastModified: now,
			Links: &[]oscal.Link{
				{
					Href: fmt.Sprintf(controlHREF, catalog.Metadata.Version, ""),
					Rel:  "canonical",
				},
			},
			OscalVersion: oscal.Version,
			Published:    &now,
			Title:        catalog.Title,
			Version:      version,
		},
	}

	catalogGroups := []oscal.Group{}

	for _, family := range catalog.ControlFamilies {
		group := oscal.Group{
			Class:    "family",
			Controls: nil,
			ID:       family.Id,
			Title:    strings.ReplaceAll(family.Description, "\n", "\\n"),
		}

		controls := []oscal.Control{}
		for _, control := range family.Controls {
			controlTitle := strings.TrimSpace(control.Title)

			newCtl := oscal.Control{
				Class: family.Id,
				ID:    control.Id,
				Title: strings.ReplaceAll(controlTitle, "\n", "\\n"),
				Parts: &[]oscal.Part{
					{
						Name:  "statement",
						ID:    fmt.Sprintf("%s_smt", control.Id),
						Prose: control.Objective,
					},
				},
				Links: &[]oscal.Link{
					{
						Href: fmt.Sprintf(controlHREF, catalog.Metadata.Version, strings.ToLower(control.Id)),
						Rel:  "canonical",
					},
				},
			}

			var subControls []oscal.Control
			for _, ar := range control.AssessmentRequirements {
				subControl := oscal.Control{
					ID:    ar.Id,
					Title: ar.Id,
					Parts: &[]oscal.Part{
						{
							Name:  "statement",
							ID:    fmt.Sprintf("%s_smt", ar.Id),
							Prose: ar.Text,
						},
					},
				}

				if ar.Recommendation != "" {
					*subControl.Parts = append(*subControl.Parts, oscal.Part{
						Name:  "guidance",
						ID:    fmt.Sprintf("%s_gdn", ar.Id),
						Prose: ar.Recommendation,
					})
				}

				*subControl.Parts = append(*subControl.Parts, oscal.Part{
					Name: "assessment-objective",
					ID:   fmt.Sprintf("%s_obj", ar.Id),
					Links: &[]oscal.Link{
						{
							Href: fmt.Sprintf("#%s_smt", ar.Id),
							Rel:  "assessment-for",
						},
					},
				})

				subControls = append(subControls, subControl)
			}

			if len(subControls) > 0 {
				newCtl.Controls = &subControls
			}
			controls = append(controls, newCtl)
		}

		group.Controls = &controls
		catalogGroups = append(catalogGroups, group)
	}
	oscalCatalog.Groups = &catalogGroups

	return oscalCatalog, nil
}
