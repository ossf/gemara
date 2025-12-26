package export

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/ossf/gemara"
	"github.com/ossf/gemara/oscal"
)

const (
	// defaultControlHrefFormat is the default URL template for linking to controls
	// in Catalog conversion. Format: controlHREF(version, controlID)
	defaultControlHrefFormat = "https://example/versions/%s#%s"
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

	profileDir := filepath.Dir(*profileOutputFile)
	catalogAbsPath, err := filepath.Abs(*catalogOutputFile)
	if err != nil {
		return fmt.Errorf("error resolving absolute path for catalog output: %w", err)
	}
	profileAbsDir, err := filepath.Abs(profileDir)
	if err != nil {
		return fmt.Errorf("error resolving absolute path for profile directory: %w", err)
	}
	relativeCatalogPath, err := filepath.Rel(profileAbsDir, catalogAbsPath)
	if err != nil {
		return fmt.Errorf("error calculating relative path: %w", err)
	}
	relativeCatalogPath = filepath.ToSlash(relativeCatalogPath)

	catalog, profile, err := oscal.FromGuidance(&guidanceDocument, relativeCatalogPath)
	if err != nil {
		return err
	}

	catalogModel := oscalTypes.OscalModels{
		Catalog: &catalog,
	}
	if err := WriteOSCALFile(catalogModel, *catalogOutputFile); err != nil {
		return err
	}

	profileModel := oscalTypes.OscalModels{
		Profile: &profile,
	}
	return WriteOSCALFile(profileModel, *profileOutputFile)
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

	oscalCatalog, err := oscal.FromCatalog(catalog, oscal.WithControlHref(defaultControlHrefFormat))
	if err != nil {
		return err
	}

	oscalModel := oscalTypes.OscalModels{
		Catalog: &oscalCatalog,
	}

	return WriteOSCALFile(oscalModel, *outputFile)
}

func WriteOSCALFile(model oscalTypes.OscalModels, outputFile string) error {
	oscalJSON, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, oscalJSON, 0600); err != nil {
		return err
	}

	fmt.Printf("Successfully wrote OSCAL content to %s\n", outputFile)
	return nil
}
