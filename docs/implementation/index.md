---
layout: page
title: Gemara Implementation
---

**Status**: <span class="badge badge-active">Active Development</span>

## Layer Schemas

Machine-readable schemas (CUE format) standardize the expression of elements in the model. Click on a layer to view its schema:

<div class="layer-grid">
  <a href="https://github.com/ossf/gemara/blob/main/schemas/layer-1.cue" class="layer-card">
    <h3>Layer 1: Guidance</h3>
    <p>High-level guidance on cybersecurity measures from industry groups and standards bodies.</p>
  </a>

  <a href="https://github.com/ossf/gemara/blob/main/schemas/layer-2.cue" class="layer-card">
    <h3>Layer 2: Controls</h3>
    <p>Technology-specific, threat-informed security controls for protecting information systems.</p>
  </a>

  <a href="https://github.com/ossf/gemara/blob/main/schemas/layer-3.cue" class="layer-card">
    <h3>Layer 3: Policy</h3>
    <p>Risk-informed guidance tailored to your organization's specific needs and risk appetite.</p>
  </a>

  <a href="https://github.com/ossf/gemara/blob/main/schemas/layer-4.cue" class="layer-card">
    <h3>Layer 4: Evaluation</h3>
    <p>Inspection of code, configurations, and deployments against policies and controls.</p>
  </a>

  <div class="layer-card">
    <h3>Layer 5: Enforcement</h3>
    <p>Prevention or remediation based on assessment findings. (Coming Soon)</p>
  </div>

  <div class="layer-card">
    <h3>Layer 6: Audit</h3>
    <p>Review of organizational policy and conformance. (Coming Soon)</p>
  </div>
</div>

**[Browse all schemas on GitHub →](https://github.com/ossf/gemara/tree/main/schemas)**

### Validation

Validate data against Gemara schemas using CUE:

```bash
go install cuelang.org/go/cmd/cue@latest
cue vet ./your-controls.yaml ./schemas/layer-2.cue
```

## Go Library

Go libraries provide APIs for reading, writing, and manipulating Gemara documents.

**[Go Package Reference →](https://pkg.go.dev/github.com/ossf/gemara)**

### Installation

```bash
go get github.com/ossf/gemara
```

### Usage Example

```go
package main

import (
    "fmt"
    "github.com/ossf/gemara"
)

func main() {
    catalog := &Catalog{}
    catalog, err := catalog.LoadFile("file://controls.yaml")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Catalog: %s\n", catalog.Metadata.ID)
}
```

See [repository examples](https://github.com/ossf/gemara/tree/main/test-data) for more.

## Contributing

The Implementation evolves based on community needs:

- **Schema improvements?** Open an issue or submit a PR
- **New features or APIs?** Propose changes via PR
- **Found a bug?** Report it
- **Significant architectural changes?** Document in an [ADR](../adr.html)

See the [Contributing Guide](https://github.com/ossf/gemara/blob/main/CONTRIBUTING.md) for details.

## Architecture Decisions

Significant implementation changes are documented in [Architecture Decision Records (ADRs)](../adr.html).

## Relationship to Other Components

### [The Model](../model)
Provides the conceptual foundation. Each schema corresponds to a layer in the model.

### [The Lexicon](../lexicon)
Informs Implementation design. Schema field names and SDK documentation use Lexicon definitions for consistency.
