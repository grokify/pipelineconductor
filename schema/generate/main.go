// Package main generates JSON Schema from Go structs.
//
//go:build ignore
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	"github.com/plexusone/pipelineconductor/pkg/model"
)

func main() {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		ExpandedStruct:             true,
	}

	schema := r.Reflect(&model.CheckResult{})
	schema.Version = "https://json-schema.org/draft/2020-12/schema"
	schema.ID = "https://github.com/plexusone/pipelineconductor/schema/check_result.schema.json"
	schema.Title = "CheckResult"
	schema.Description = "Workflow compliance check result"

	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling schema: %v\n", err)
		os.Exit(1)
	}

	// Determine output path
	outputPath := "schema/check_result.schema.json"
	if len(os.Args) > 1 {
		outputPath = os.Args[1]
	}

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}

	// Write schema file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Schema written to %s\n", outputPath)
}
