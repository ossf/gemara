package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ossf/gemara"
	"github.com/ossf/gemara/checklist"
)

func main() {
	policyFile := flag.String("policy", "", "Path to the policy file (YAML or JSON)")
	outputFile := flag.String("output", "-", "Output file path (default: stdout)")
	flag.Parse()

	if *policyFile == "" {
		fmt.Fprintf(os.Stderr, "No policy file supplied.\n")
		os.Exit(1)
	}

	// Load policy
	policy := &gemara.Policy{}
	policyPath := ensureFileScheme(*policyFile)
	if err := policy.LoadFile(policyPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading policy file: %v\n", err)
		os.Exit(1)
	}

	// Generate markdown (catalogs are loaded on-demand from metadata.mapping-references)
	markdown, err := checklist.ToMarkdownChecklist(*policy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating markdown: %v\n", err)
		os.Exit(1)
	}

	// Output
	if *outputFile != "-" {
		if err := os.WriteFile(*outputFile, []byte(markdown), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Checklist written to %s\n", *outputFile)
	} else {
		fmt.Print(markdown)
	}
}

// ensureFileScheme ensures the path has a file:// scheme prefix
func ensureFileScheme(path string) string {
	if len(path) >= 7 && path[:7] == "file://" {
		return path
	}
	// Convert to absolute path for file:// scheme
	absPath, err := filepath.Abs(path)
	if err != nil {
		// If we can't get absolute path, just use the original
		return "file://" + path
	}
	return "file://" + absPath
}
