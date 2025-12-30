package oscal

import (
	"testing"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/ossf/gemara"
	"github.com/stretchr/testify/assert"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

var testCases = []struct {
	name          string
	catalog       *gemara.Catalog
	controlHREF   string
	wantErr       bool
	expectedTitle string
}{
	{
		name: "Valid catalog with single control family",
		catalog: &gemara.Catalog{
			Metadata: gemara.Metadata{
				Id:      "test-catalog",
				Version: "devel",
			},
			Title: "Test Catalog",
			Families: []gemara.Family{
				{
					Id:          "AC",
					Title:       "access-control",
					Description: "Controls for access management",
				},
			},
			Controls: []gemara.Control{
				{
					Id:     "AC-01",
					Family: "AC",
					Title:  "Access Control Policy",
					AssessmentRequirements: []gemara.AssessmentRequirement{
						{
							Id:   "AC-01.1",
							Text: "Develop and document access control policy",
						},
					},
				},
			},
		},
		controlHREF:   "https://baseline.openssf.org/versions/%s#%s",
		wantErr:       false,
		expectedTitle: "Test Catalog",
	},
	{
		name: "Valid catalog with multiple control families",
		catalog: &gemara.Catalog{
			Metadata: gemara.Metadata{
				Id:      "test-catalog-multi",
				Version: "devel",
			},
			Title: "Test Catalog Multiple",
			Families: []gemara.Family{
				{
					Id:          "AC",
					Title:       "access-control",
					Description: "Controls for access management",
				},
				{
					Id:          "BR",
					Title:       "business-requirements",
					Description: "Controls for business requirements",
				},
			},
			Controls: []gemara.Control{
				{
					Id:     "AC-01",
					Family: "AC",
					Title:  "Access Control Policy",
					AssessmentRequirements: []gemara.AssessmentRequirement{
						{
							Id:   "AC-01.1",
							Text: "Develop and document access control policy",
						},
					},
				},
				{
					Id:     "BR-01",
					Family: "BR",
					Title:  "Business Requirements Policy",
					AssessmentRequirements: []gemara.AssessmentRequirement{
						{
							Id:   "BR-01.1",
							Text: "Define business requirements",
						},
					},
				},
			},
		},
		controlHREF:   "https://baseline.openssf.org/versions/%s#%s",
		wantErr:       false,
		expectedTitle: "Test Catalog Multiple",
	},
}

func TestFromCatalog(t *testing.T) {
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			oscalCatalog, err := FromCatalog(tt.catalog, WithControlHref(tt.controlHREF))

			if (err == nil) == tt.wantErr {
				t.Errorf("ToOSCAL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Wrap oscal catalog
			// Create the proper OSCAL document structure
			oscalDocument := oscalTypes.OscalModels{
				Catalog: &oscalCatalog,
			}

			// Create validation for the OSCAL catalog
			assert.NoError(t, oscalUtils.Validate(oscalDocument))

			// Compare each field
			assert.NotEmpty(t, oscalCatalog.UUID)
			assert.Equal(t, tt.expectedTitle, oscalCatalog.Metadata.Title)
			assert.Equal(t, tt.catalog.Metadata.Version, oscalCatalog.Metadata.Version)
			assert.Equal(t, len(tt.catalog.Families), len(*oscalCatalog.Groups))

			// Compare each control family
			for i, family := range tt.catalog.Families {
				groups := (*oscalCatalog.Groups)
				group := groups[i]
				assert.Equal(t, family.Id, group.ID)
			}
		})
	}
}
