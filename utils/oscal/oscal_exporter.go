package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	oscal "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
	"github.com/goccy/go-yaml"

	"github.com/ossf/gemara/layer1"
	"github.com/ossf/gemara/layer2"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 2 {
		fmt.Println("Usage: oscal_exporter <subcommand> <path> [flags]")
		fmt.Println("Available subcommands: guidance, catalog")
		os.Exit(1)
	}

	subcommand, path := args[0], args[1]
	subcommandArgs := args[2:]

	var err error
	switch subcommand {
	case "guidance":
		err = exportGuidance(path, subcommandArgs)
	case "catalog":
		err = exportCatalog(path, subcommandArgs)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		os.Exit(1)
	}
}

func exportGuidance(path string, args []string) error {
	cmd := flag.NewFlagSet("guidance", flag.ExitOnError)
	outputFile := cmd.String("output", "guidance.json", "Path to output file for OSCAL Catalog and Profile")
	if err := cmd.Parse(args); err != nil {
		return err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var guidanceDocument layer1.GuidanceDocument
	if err := yaml.Unmarshal(data, &guidanceDocument); err != nil {
		return err
	}

	oscalCatalog, err := guidanceDocument.ToOSCALCatalog()
	if err != nil {
		return err
	}

	oscalProfile, err := guidanceDocument.ToOSCALProfile(fmt.Sprintf("file://%s", *outputFile))
	if err != nil {
		return err
	}

	oscalModel := oscal.OscalModels{
		Catalog: &oscalCatalog,
		Profile: &oscalProfile,
	}

	return writeOSCALFile(oscalModel, *outputFile)
}

func exportCatalog(path string, args []string) error {
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
