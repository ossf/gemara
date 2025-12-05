package oscal

import "github.com/ossf/gemara"

type generateOpts struct {
	version       string
	imports       map[string]string
	canonicalHref string
	controlHREF   string
}

func (g *generateOpts) completeFromGuidance(doc gemara.GuidanceDocument) {
	if g.version == "" {
		g.version = doc.Metadata.Version
	}
	if g.imports == nil {
		g.imports = make(map[string]string)
		for _, mappingRef := range doc.Metadata.MappingReferences {
			g.imports[mappingRef.Id] = mappingRef.Url
		}
	}
}

func (g *generateOpts) completeFromCatalog(catalog *gemara.Catalog) {
	if g.version == "" {
		g.version = catalog.Metadata.Version
	}
}

// GenerateOption defines an option to tune the behavior of the OSCAL
// generation functions for both Layer 1 (GuidanceDocument) and Layer 2 (Catalog).
type GenerateOption func(opts *generateOpts)

// WithVersion is a GenerateOption that sets the version of the OSCAL Document. If set,
// this will be used instead of the version in GuidanceDocument.
func WithVersion(version string) GenerateOption {
	return func(opts *generateOpts) {
		opts.version = version
	}
}

// WithOSCALImports is a GenerateOption that provides the `href` to guidance document mappings in OSCAL
// by mapping unique identifier. If unset, the mapping URL of the guidance document will be used.
func WithOSCALImports(imports map[string]string) GenerateOption {
	return func(opts *generateOpts) {
		opts.imports = imports
	}
}

// WithCanonicalHrefFormat is a GenerateOption that provides an `href` format string
// for the canonical version of the guidance document. If set, this will be added as a
// link in the mapping.cue with the rel="canonical" attribute. Ex - https://myguidance.org/versions/%s
func WithCanonicalHrefFormat(canonicalHref string) GenerateOption {
	return func(opts *generateOpts) {
		opts.canonicalHref = canonicalHref
	}
}

// WithControlHref is a GenerateOption that provides a URL template for linking to controls
// in Catalog conversion. Uses format: controlHREF(version, controlID)
// Example: "https://baseline.openssf.org/versions/%s#%s"
func WithControlHref(controlHref string) GenerateOption {
	return func(opts *generateOpts) {
		opts.controlHREF = controlHref
	}
}
