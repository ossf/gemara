---
layout: page
title: Gemara Implementation
---

**Status**: <span class="badge badge-active">Active Development</span>

## Overview

The Implementation provides the practical tools for working with Gemara:
- **Schemas**: Machine-readable formats (CUE) for validation and standardization
- **Go Library**: Go library for working with Gemara documents programmatically

## Components

### Schemas

Machine-readable schemas (CUE format) that standardize the expression of different elements 
in the implemented model.

**Available Schemas:**
- Layer 1 Schema (`schemas/layer-1.cue`) - Guidance documents
- Layer 2 Schema (`schemas/layer-2.cue`) - Control catalogs
- Layer 3 Schema (`schemas/layer-3.cue`) - Policy documents
- Layer 4 Schema (`schemas/layer-4.cue`) - Evaluation results

**[Browse Schemas on GitHub →](https://github.com/ossf/gemara/tree/main/schemas)**

Validate your data against Gemara schemas using CUE:

```bash
# Install CUE
go install cuelang.org/go/cmd/cue@latest

# Validate a Layer 2 control catalog
cue vet ./your-controls.yaml ./schemas/layer-2.cue
```

### Go Library

Programming libraries that provide APIs for working with Gemara documents with support for reading, writing, and manipulating Gemara documents.

**[Go Package Reference →](https://pkg.go.dev/github.com/ossf/gemara)**

Install the Go module:

```bash
go get github.com/ossf/gemara
```

Load and work with documents:

```go
package main

import (
    "fmt"
    "github.com/ossf/gemara"
)

func main() {
    // Load a control catalog 
    catalog := &Catalog{}
    catalog, err := catalog.LoadFile("file://controls.yaml")
    if err != nil {
        panic(err)
    }
    
    // Work with the catalog 
    fmt.Printf("Catalog: %s\n", catalog.Metadata.ID)   
}
```

For more examples, see the [repository examples](https://github.com/ossf/gemara/tree/main/test-data).

## Relationship to Other Components

As one of Gemara's three core components, the Implementation works alongside:

### [The Model Component](../model)
The Model component provides the conceptual foundation. Each schema 
corresponds to a layer in the model, and the SDK reflects the model's structure.

### [The Lexicon Component](../lexicon)
The Lexicon component informs the Implementation's design. Schema field 
names and SDK documentation use Lexicon definitions to ensure consistency and shared 
understanding.

Together, these three components support the entire Gemara ecosystem.
