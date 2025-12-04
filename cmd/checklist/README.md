# Checklist Generator

A command-line tool to generate evaluation checklists from Gemara Policy documents.

## Usage

```bash
checklist -policy <policy-file> [-output <output-file>]
```

## Options

- `-policy` (required): Path to the policy file (YAML or JSON)
- `-output` (optional): Output file path (default: stdout)

## Examples

### Basic usage

```bash
checklist -policy ./policy.yaml
```

### Output to file

```bash
checklist -policy ./policy.yaml -output evaluation-checklist.md
```

## How it works

1. Loads the policy file specified by `-policy`
2. Generates a checklist from the policy's assessment requirement modifications
3. For each assessment requirement:
   - Uses override text if available
   - Otherwise loads description from Layer 2 catalog on-demand (from `metadata.mapping-references` if recommendation not overridden)
   - Catalogs are cached in memory for reuse during the operation
   - Includes all required evaluators as checklist items
4. Outputs the checklist as Markdown

## Catalog Loading

Catalogs are automatically loaded on-demand from the `metadata.mapping-references` section of the policy file. Each `MappingReference` should have:
- `id`: The reference ID (e.g., "OSPS-B")
- `url`: Path or URL to the catalog file (supports `file://` paths and `https://` URLs)

Catalogs are only loaded when needed (when a requirement doesn't have an override) and are cached for the duration of the operation to avoid redundant file I/O.

## Output Format

The generated checklist includes:
- Policy ID and author information
- Control sections organized by reference ID
- Requirement subsections with recommendations/descriptions
- Evaluator items as checkboxes
