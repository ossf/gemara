package export

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/ossf/gemara/layer1"
	"github.com/ossf/gemara/layer2"
)

func Guidance(path string, args []string) error {
	cmd := flag.NewFlagSet("guidance", flag.ExitOnError)
	catalogOutputFile := cmd.String("catalog-output", "guidance.json", "Path to output file for OSCAL Catalog")
	profileOutputFile := cmd.String("profile-output", "profile.json", "Path to output file for OSCAL Profile")
	if err := cmd.Parse(args); err != nil {
		return err
	}

	var guidanceDocument layer1.GuidanceDocument
	pathWithScheme := fmt.Sprintf("file://%s", path)
	if err := guidanceDocument.LoadFile(pathWithScheme); err != nil {
		return err
	}

	oscalCatalog, err := guidanceDocument.ToOSCALCatalog()
	if err != nil {
		return err
	}

	oscalProfile, err := guidanceDocument.ToOSCALProfile(fmt.Sprintf("file://%s", *catalogOutputFile))
	if err != nil {
		return err
	}

	catalogOscalModel := oscal.OscalModels{
		Catalog: &oscalCatalog,
	}

	if err := writeOSCALFile(catalogOscalModel, *catalogOutputFile); err != nil {
		return err
	}

	profileOscalModel := oscal.OscalModels{
		Profile: &oscalProfile,
	}

	return writeOSCALFile(profileOscalModel, *profileOutputFile)
}

func Catalog(path string, args []string) error {
	cmd := flag.NewFlagSet("catalog", flag.ExitOnError)
	outputFile := cmd.String("output", "catalog.json", "Path to output file")
	if err := cmd.Parse(args); err != nil {
		return err
	}

	catalog := &layer2.Catalog{}
	pathWithScheme := fmt.Sprintf("file://%s", path)
	if err := catalog.LoadFile(pathWithScheme); err != nil {
		return err
	}

	oscalCatalog, err := catalog.ToOSCAL("https://example/versions/%s#%s")
	if err != nil {
		return err
	}

	oscalModel := oscal.OscalModels{
		Catalog: &oscalCatalog,
	}

	return writeOSCALFile(oscalModel, *outputFile)
}

func writeOSCALFile(model oscal.OscalModels, outputFile string) error {
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
