package export

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/ossf/gemara"
	"github.com/ossf/gemara/oscal"
)

func Guidance(path string, args []string) error {
	cmd := flag.NewFlagSet("guidance", flag.ExitOnError)
	catalogOutputFile := cmd.String("catalog-output", "guidance.json", "Path to output file for OSCAL ExportCatalog")
	profileOutputFile := cmd.String("profile-output", "profile.json", "Path to output file for OSCAL Profile")
	if err := cmd.Parse(args); err != nil {
		return err
	}

	var guidanceDocument gemara.GuidanceDocument
	pathWithScheme := fmt.Sprintf("file://%s", path)
	if err := guidanceDocument.LoadFile(pathWithScheme); err != nil {
		return err
	}

	oscalCatalog, err := oscal.CatalogFromGuidanceDocument(&guidanceDocument)
	if err != nil {
		return err
	}

	oscalProfile, err := oscal.ProfileFromGuidanceDocument(&guidanceDocument, fmt.Sprintf("file://%s", *catalogOutputFile))
	if err != nil {
		return err
	}

	catalogOscalModel := oscalTypes.OscalModels{
		Catalog: &oscalCatalog,
	}

	if err := WriteOSCALFile(catalogOscalModel, *catalogOutputFile); err != nil {
		return err
	}

	profileOscalModel := oscalTypes.OscalModels{
		Profile: &oscalProfile,
	}

	return WriteOSCALFile(profileOscalModel, *profileOutputFile)
}

func Catalog(path string, args []string) error {
	cmd := flag.NewFlagSet("catalog", flag.ExitOnError)
	outputFile := cmd.String("output", "catalog.json", "Path to output file")
	if err := cmd.Parse(args); err != nil {
		return err
	}

	catalog := &gemara.Catalog{}
	pathWithScheme := fmt.Sprintf("file://%s", path)
	if err := catalog.LoadFile(pathWithScheme); err != nil {
		return err
	}

	oscalCatalog, err := oscal.FromCatalog(catalog, "https://example/versions/%s#%s")
	if err != nil {
		return err
	}

	oscalModel := oscalTypes.OscalModels{
		Catalog: &oscalCatalog,
	}

	return WriteOSCALFile(oscalModel, *outputFile)
}

func WriteOSCALFile(model oscalTypes.OscalModels, outputFile string) error {
	oscalJSON, err := json.MarshalIndent(model, "", "  ") // Using " " for indent
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, oscalJSON, 0600); err != nil {
		return err
	}

	fmt.Printf("Successfully wrote OSCAL content to %s\n", outputFile)
	return nil
}
