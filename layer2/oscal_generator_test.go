package layer2

import (
	"testing"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/stretchr/testify/assert"

	oscalUtils "github.com/ossf/gemara/internal/oscal"
)

var TestCases = []struct {
	name          string
	catalog       *Catalog
	controlHREF   string
	wantErr       bool
	expectedTitle string
}{
	{
		name: "Valid catalog with single control family",
		catalog: &Catalog{
			Metadata: Metadata{
				Id:      "test-catalog",
				Title:   "Test Catalog",
				Version: "devel",
			},
			ControlFamilies: []ControlFamily{
				{
					Id:          "AC",
					Title:       "access-control",
					Description: "Controls for access management",
					Controls: []Control{
						{
							Id:    "AC-01",
							Title: "Access Control Policy",
							AssessmentRequirements: []AssessmentRequirement{
								{
									Id:   "AC-01.1",
									Text: "Develop and document access control policy",
								},
							},
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
		catalog: &Catalog{
			Metadata: Metadata{
				Id:      "test-catalog-multi",
				Title:   "Test Catalog Multiple",
				Version: "devel",
			},
			ControlFamilies: []ControlFamily{
				{
					Id:          "AC",
					Title:       "access-control",
					Description: "Controls for access management",
					Controls: []Control{
						{
							Id:    "AC-01",
							Title: "Access Control Policy",
							AssessmentRequirements: []AssessmentRequirement{
								{
									Id:   "AC-01.1",
									Text: "Develop and document access control policy",
								},
							},
						},
					},
				},
				{
					Id:          "BR",
					Title:       "business-requirements",
					Description: "Controls for business requirements",
					Controls: []Control{
						{
							Id:    "BR-01",
							Title: "Business Requirements Policy",
							AssessmentRequirements: []AssessmentRequirement{
								{
									Id:   "BR-01.1",
									Text: "Define business requirements",
								},
							},
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

func Test_toOSCAL(t *testing.T) {
	for _, tt := range TestCases {
		t.Run(tt.name, func(t *testing.T) {
			oscalCatalog, err := tt.catalog.ToOSCAL(tt.controlHREF)

			if (err == nil) == tt.wantErr {
				t.Errorf("ToOSCAL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Wrap oscal catalog
			// Create the proper OSCAL document structure
			oscalDocument := oscal.OscalModels{
				Catalog: &oscalCatalog,
			}

			// Create validation for the OSCAL catalog
			assert.NoError(t, oscalUtils.Validate(oscalDocument))

			// Compare each field
			assert.NotEmpty(t, oscalCatalog.UUID)
			assert.Equal(t, tt.expectedTitle, oscalCatalog.Metadata.Title)
			assert.Equal(t, tt.catalog.Metadata.Version, oscalCatalog.Metadata.Version)
			assert.Equal(t, len(tt.catalog.ControlFamilies), len(*oscalCatalog.Groups))

			// Compare each control family
			for i, family := range tt.catalog.ControlFamilies {
				groups := (*oscalCatalog.Groups)
				group := groups[i]
				assert.Equal(t, family.Id, group.ID) 
			}
		})
	}
}
