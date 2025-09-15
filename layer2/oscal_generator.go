package layer2

import (
	"fmt"
	"strings"
	"time"

	"github.com/defenseunicorns/go-oscal/src/pkg/uuid"
	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

// ToOSCAL converts a Catalog to OSCAL Catalog format.
// Parameters:
//   - controlHREF: URL template for linking to controls. Uses format: controlHREF(version, controlID)
//     Example: "https://baseline.openssf.org/versions/%s#%s"
//
// The function automatically:
//   - Uses the catalog's internal version from Metadata.Version
//   - Uses the ControlFamily.Id as the OSCAL group ID
//   - Generates a unique UUID for the catalog
func (c *Catalog) ToOSCAL(controlHREF string) (oscal.Catalog, error) {
	now := time.Now()

	oscalCatalog := oscal.Catalog{
		UUID:   uuid.NewUUID(),
		Groups: nil,
		Metadata: oscal.Metadata{
			LastModified: oscalUtils.GetTimeWithFallback(c.Metadata.LastModified, now),
			Links: &[]oscal.Link{
				{
					Href: fmt.Sprintf(controlHREF, c.Metadata.Version, ""),
					Rel:  "canonical",
				},
			},
			OscalVersion: oscal.Version,
			Published:    &now,
			Title:        c.Metadata.Title,
			Version:      c.Metadata.Version,
		},
	}

	catalogGroups := []oscal.Group{}

	for _, family := range c.ControlFamilies {
		group := oscal.Group{
			Class:    "family",
			Controls: nil,
			ID:       family.Id,
			Title:    family.Description,
		}

		controls := []oscal.Control{}
		for _, control := range family.Controls {
			parts := []oscal.Part{}
			for _, ar := range control.AssessmentRequirements {
				parts = append(parts, oscal.Part{
					Class: control.Id,
					ID:    ar.Id,
					Name:  ar.Id,
					Ns:    "",
					Parts: &[]oscal.Part{
						{
							ID:    ar.Id + ".R",
							Name:  "recommendation",
							Ns:    oscalUtils.GemaraNamespace,
							Prose: ar.Recommendation,
							Links: &[]oscal.Link{
								{
									Href: fmt.Sprintf(controlHREF, c.Metadata.Version, ar.Id),
									Rel:  "canonical",
								},
							},
						},
					},
					Prose: ar.Text,
					Title: "",
				})
			}

			newCtl := oscal.Control{
				Class: family.Title,
				ID:    control.Id,
				Links: &[]oscal.Link{
					{
						Href: fmt.Sprintf(controlHREF, c.Metadata.Version, strings.ToLower(control.Id)),
						Rel:  "canonical",
					},
				},
				Parts: &parts,
				Title: strings.TrimSpace(control.Title),
			}
			controls = append(controls, newCtl)
		}

		group.Controls = &controls
		catalogGroups = append(catalogGroups, group)
	}
	oscalCatalog.Groups = &catalogGroups

	return oscalCatalog, nil
}
