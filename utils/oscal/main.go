package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ossf/gemara/utils/oscal/export"
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
		err = export.Guidance(path, subcommandArgs)
	case "catalog":
		err = export.Catalog(path, subcommandArgs)
	default:
		fmt.Printf("Unknown subcommand: %s\n", subcommand)
		os.Exit(1)
	}

	if err != nil {
		fmt.Printf("Error processing command: %v\n", err)
		os.Exit(1)
	}
}
