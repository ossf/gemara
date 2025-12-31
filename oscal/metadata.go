package oscal

import (
	"fmt"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/ossf/gemara"
	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

// createMetadata creates OSCAL metadata with common fields and optional author information
func createMetadata(title string, version string, published *time.Time, canonicalHref string, authorName string) oscal.Metadata {
	now := time.Now()
	// Ensure version is never empty by using default if not provided
	if version == "" {
		version = oscalUtils.DefaultOSCALVersion
	}
	metadata := oscal.Metadata{
		Title:        title,
		OscalVersion: oscal.Version,
		Version:      version,
		Published:    published,
		LastModified: now,
	}

	if canonicalHref != "" {
		metadata.Links = &[]oscal.Link{
			{
				Href: canonicalHref,
				Rel:  "canonical",
			},
		}
	}

	if authorName != "" {
		authorRole := oscal.Role{
			ID:          "author",
			Description: "Author and owner of the document",
			Title:       "Author",
		}

		author := oscal.Party{
			UUID: uuid.NewUUID(),
			Type: "person",
			Name: authorName,
		}

		responsibleParty := oscal.ResponsibleParty{
			PartyUuids: []string{author.UUID},
			RoleId:     authorRole.ID,
		}

		metadata.Parties = &[]oscal.Party{author}
		metadata.Roles = &[]oscal.Role{authorRole}
		metadata.ResponsibleParties = &[]oscal.ResponsibleParty{responsibleParty}
	}

	return metadata
}

// createMetadataFromGuidance creates OSCAL metadata from a GuidanceDocument
func createMetadataFromGuidance(guidance *gemara.GuidanceDocument, opts generateOpts) (oscal.Metadata, error) {
	canonicalHref := ""
	if opts.canonicalHref != "" {
		canonicalHref = fmt.Sprintf(opts.canonicalHref, opts.version)
	}

	published := oscalUtils.GetTime(string(guidance.Metadata.Date))
	metadata := createMetadata(
		guidance.Title,
		opts.version,
		published,
		canonicalHref,
		guidance.Metadata.Author.Name,
	)

	return metadata, nil
}

// createMetadataFromCatalog creates OSCAL metadata from a Catalog
func createMetadataFromCatalog(catalog *gemara.Catalog, opts generateOpts) (oscal.Metadata, error) {
	// Handle canonical HREF - prefer controlHREF for Catalog, fallback to canonicalHref
	var canonicalHref string
	if opts.controlHREF != "" {
		// controlHREF format: "https://example.com/versions/%s#%s" (version, controlID)
		canonicalHref = fmt.Sprintf(opts.controlHREF, opts.version, "")
	} else if opts.canonicalHref != "" {
		// canonicalHref format: "https://example.com/versions/%s" (version only)
		canonicalHref = fmt.Sprintf(opts.canonicalHref, opts.version)
	}

	published := oscalUtils.GetTime(string(catalog.Metadata.Date))
	metadata := createMetadata(
		catalog.Title,
		opts.version,
		published,
		canonicalHref,
		catalog.Metadata.Author.Name,
	)

	return metadata, nil
}
